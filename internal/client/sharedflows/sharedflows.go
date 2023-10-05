// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sharedflows

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"internal/apiclient"

	"internal/clilog"
)

type sharedflows struct {
	Flows []sharedflow `json:"sharedFlows,omitempty"`
}

type sharedflow struct {
	Name     string   `json:"name,omitempty"`
	Revision []string `json:"revision,omitempty"`
}

type revision struct {
	name string
	rev  string
}

// Create
func Create(name string, proxy string) (respBody []byte, err error) {
	if proxy != "" {
		err = apiclient.ImportBundle("sharedflows", name, proxy)
		return respBody, err
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows")
	proxyName := "{\"name\":\"" + name + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), proxyName)
	return respBody, err
}

// Get
func Get(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if revision != -1 {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name, "revisions", strconv.Itoa(revision))
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name)
	}
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetHighestSfRevision
func GetHighestSfRevision(name string) (version int, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name)
	respBody, err := apiclient.HttpClient(u.String())
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	if err != nil {
		return -1, err
	}

	sfRevisions := sharedflow{}
	if err = json.Unmarshal(respBody, &sfRevisions); err != nil {
		return -1, err
	}
	version, err = strconv.Atoi(maxRevision(sfRevisions.Revision))
	if err != nil {
		return -1, nil
	}
	return version, nil
}

// Delete
func Delete(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if revision != -1 {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name, "revisions", strconv.Itoa(revision))
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name)
	}
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List(includeRevisions bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if includeRevisions {
		q := u.Query()
		q.Set("includeRevisions", strconv.FormatBool(includeRevisions))
		u.RawQuery = q.Encode()
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListEnvDeployments
func ListEnvDeployments() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}

	q := u.Query()
	q.Set("sharedFlows", "true")
	u.RawQuery = q.Encode()

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "deployments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListDeployments
func ListDeployments(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name, "deployments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListRevisionDeployments
func ListRevisionDeployments(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "sharedflows", name, "revisions",
		strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Deploy
func Deploy(name string, revision int, overrides bool, serviceAccountName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if overrides || serviceAccountName != "" {
		q := u.Query()
		if overrides {
			q.Set("override", "true")
		}
		if serviceAccountName != "" {
			q.Set("serviceAccount", serviceAccountName)
		}
		u.RawQuery = q.Encode()
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"sharedflows", name, "revisions", strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

// Clean
func Clean(name string, reportOnly bool) (err error) {
	type deploymentsDef struct {
		Environment    string `json:"environment,omitempty"`
		APIProxy       string `json:"apiProxy,omitempty"`
		Revision       string `json:"revision,omitempty"`
		DeloyStartTime string `json:"deployStartTime,omitempty"`
		BasePath       string `json:"basePath,omitempty"`
	}

	type sfDeploymentsDef struct {
		Deployments []deploymentsDef `json:"deployments,omitempty"`
	}

	type metaDataDef struct {
		CreatedAt      string `json:"createdAt,omitempty"`
		LastModifiedAt string `json:"lastModifiedAt,omitempty"`
		SubType        string `json:"subType,omitempty"`
	}

	type sfRevisionsDef struct {
		MetaData metaDataDef `json:"metaData,omitempty"`
		Name     string      `json:"name,omitempty"`
		Revision []string    `json:"revision,omitempty"`
	}

	sfDeployments := sfDeploymentsDef{}
	sfRevisions := sfRevisionsDef{}

	reportRevisions := make(map[string]bool)
	deployedRevisions := make(map[string]bool)

	var sfDeploymentsBytes, sfRevisionsBytes []byte
	var revision int

	// disable printing
	apiclient.ClientPrintHttpResponse.Set(false)

	// step 1. get a list of revisions that are deployed.
	if sfDeploymentsBytes, err = ListDeployments(name); err != nil {
		return err
	}

	if err = json.Unmarshal(sfDeploymentsBytes, &sfDeployments); err != nil {
		return err
	}

	if len(sfDeployments.Deployments) == 0 {
		return fmt.Errorf("no sharedflow deployments found")
	}

	for _, sfDeployment := range sfDeployments.Deployments {
		if !deployedRevisions[sfDeployment.Revision] {
			deployedRevisions[sfDeployment.Revision] = true
		}
	}

	clilog.Info.Println("Revisions [" + getRevisions(deployedRevisions) + "] deployed for Sharedflow " + name)

	// step 2. get all the revisions for the sf
	if sfRevisionsBytes, err = Get(name, -1); err != nil {
		return err
	}

	if err = json.Unmarshal(sfRevisionsBytes, &sfRevisions); err != nil {
		return err
	}

	// enable printing
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	for _, sfRevision := range sfRevisions.Revision {
		if !isRevisionDeployed(deployedRevisions, sfRevision) {
			// step 3. clean up proxy revisions that are not deployed
			if reportOnly {
				if !reportRevisions[sfRevision] {
					reportRevisions[sfRevision] = true
				}
			} else {
				if revision, err = strconv.Atoi(sfRevision); err != nil {
					return err
				}
				clilog.Info.Println("Deleting revision: " + sfRevision)
				if _, err = Delete(name, revision); err != nil {
					return err
				}
			}
		}
	}

	if reportOnly && len(reportRevisions) > 0 {
		clilog.Info.Println("[REPORT]: Sharedflow '" + name + "' revisions: " + getRevisions(reportRevisions) + " can be cleaned")
	}

	return nil
}

// Undeploy
func Undeploy(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"sharedflows", name, "revisions", strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Fetch
func Fetch(name string, revision int) (err error) {
	return apiclient.FetchBundle("sharedflows", "", name, strconv.Itoa(revision), true)
}

// Export
func Export(conn int, folder string, allRevisions bool) (err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	q := u.Query()
	q.Set("includeRevisions", "true")
	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows")

	// don't print to sysout
	apiclient.ClientPrintHttpResponse.Set(false)
	respBody, err := apiclient.HttpClient(u.String())
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	if err != nil {
		return err
	}

	var shrdflows sharedflows
	if err = json.Unmarshal(respBody, &shrdflows); err != nil {
		return err
	}

	clilog.Debug.Printf("Found %d API Proxies in the org\n", len(shrdflows.Flows))
	clilog.Debug.Printf("Exporting bundles with %d connections\n", conn)

	jobChan := make(chan revision)
	errChan := make(chan error)

	fanOutWg := sync.WaitGroup{}
	fanInWg := sync.WaitGroup{}

	errs := []string{}
	fanInWg.Add(1)
	go func() {
		defer fanInWg.Done()
		for {
			newErr, ok := <-errChan
			if !ok {
				return
			}
			errs = append(errs, newErr.Error())
		}
	}()

	for i := 0; i < conn; i++ {
		fanOutWg.Add(1)
		go exportSharedFlows(&fanOutWg, jobChan, folder, allRevisions, errChan)
	}

	for _, proxy := range shrdflows.Flows {
		if allRevisions {
			for _, rev := range proxy.Revision {
				jobChan <- revision{name: proxy.Name, rev: rev}
			}
		} else {
			lastRevision := maxRevision(proxy.Revision)
			jobChan <- revision{name: proxy.Name, rev: lastRevision}
		}
	}
	close(jobChan)
	fanOutWg.Wait()
	close(errChan)
	fanInWg.Wait()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func exportSharedFlows(wg *sync.WaitGroup, jobs <-chan revision, folder string, allRevisions bool, errs chan<- error) {
	defer wg.Done()
	for {
		job, ok := <-jobs
		if !ok {
			return
		}
		err := apiclient.FetchBundle("sharedflows", folder, job.name, job.rev, allRevisions)
		if err != nil {
			errs <- err
		}
	}
}

// Import
func Import(conn int, folder string) error {
	var bundles []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			clilog.Warning.Println("sharedflows folder not found")
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".zip" {
			return nil
		}
		bundles = append(bundles, path)
		return nil
	})
	if err != nil {
		return err
	}

	clilog.Debug.Printf("Found %d sharedflow bundles in the folder\n", len(bundles))
	clilog.Debug.Printf("Importing sharedflows with %d parallel connections\n", conn)

	jobChan := make(chan string)
	errChan := make(chan error)

	fanOutWg := sync.WaitGroup{}
	fanInWg := sync.WaitGroup{}

	errs := []string{}
	fanInWg.Add(1)
	go func() {
		defer fanInWg.Done()
		for {
			newErr, ok := <-errChan
			if !ok {
				return
			}
			errs = append(errs, newErr.Error())
		}
	}()

	for i := 0; i < conn; i++ {
		fanOutWg.Add(1)
		go importSharedFlows(&fanOutWg, jobChan, errChan)
	}

	for _, bundle := range bundles {
		jobChan <- bundle
	}
	close(jobChan)
	fanOutWg.Wait()
	close(errChan)
	fanInWg.Wait()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func importSharedFlows(wg *sync.WaitGroup, jobs <-chan string, errs chan<- error) {
	defer wg.Done()
	for {
		job, ok := <-jobs
		if !ok {
			return
		}

		u, _ := url.Parse(apiclient.BaseURL)
		q := u.Query()
		q.Set("name", strings.TrimSuffix(filepath.Base(job), ".zip"))
		q.Set("action", "import")
		u.RawQuery = q.Encode()
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows")

		fd, err := os.OpenFile(job, os.O_RDONLY|os.O_EXCL, 0o644)
		if err != nil {
			errs <- err
			continue
		}
		reqBody := &bytes.Buffer{}
		w := multipart.NewWriter(reqBody)
		part, err := w.CreateFormFile("file", filepath.Base(job))
		if err != nil {
			errs <- err
			continue
		}
		_, _ = io.Copy(part, fd)
		fd.Close()
		w.Close()

		err = apiclient.GetHttpClient()
		if err != nil {
			errs <- err
			continue
		}
		req, err := http.NewRequest(http.MethodPost, u.String(), reqBody)
		if err != nil {
			errs <- err
			continue
		}
		req, err = apiclient.SetAuthHeader(req)
		if err != nil {
			errs <- err
			continue
		}
		req.Header.Add("Content-Type", w.FormDataContentType())

		resp, err := apiclient.ApigeeAPIClient.Do(req)
		if err != nil {
			errs <- err
			continue
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			b = []byte(err.Error())
		}
		if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
			errs <- fmt.Errorf("bundle not imported: (HTTP %v) %s", resp.StatusCode, b)
			continue
		}

		if len(b) > 0 && apiclient.GetPrintOutput() {
			out := bytes.NewBuffer([]byte{})
			if err = json.Indent(out, bytes.TrimSpace(b), "", "  "); err != nil {
				errs <- fmt.Errorf("apigee returned invalid json: %w", err)
			}
		}
		clilog.Debug.Printf("Completed bundle import: %s", job)
	}
}

func isRevisionDeployed(revisions map[string]bool, revision string) bool {
	for r := range revisions {
		if r == revision {
			return true
		}
	}
	return false
}

func getRevisions(r map[string]bool) string {
	var arr []string
	for s := range r {
		arr = append(arr, s)
	}
	return strings.Join(arr, ",")
}

func maxRevision(revisionList []string) string {
	max := 1
	for i := 0; i < len(revisionList); i++ {
		revisionInt, _ := strconv.Atoi(revisionList[i])
		if max < revisionInt {
			max = revisionInt
		}
	}
	return strconv.Itoa(max)
}

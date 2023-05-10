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

package apis

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

type proxies struct {
	Proxies []proxy `json:"proxies,omitempty"`
}

type proxy struct {
	Name     string   `json:"name,omitempty"`
	Revision []string `json:"revision,omitempty"`
}

type revision struct {
	name string
	rev  string
}

type deploychangereport struct {
	ValidationErrors validationerrors   `json:"validationErrors,omitempty"`
	RoutingChanges   []routingchanges   `json:"routingChanges,omitempty"`
	RoutingConflicts []routingconflicts `json:"routingConflicts,omitempty"`
}

type validationerrors struct {
	Violations []violation `json:"violation,omitempty"`
}

type violation struct {
	Type        string `json:"type,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
}

type routingchanges struct {
}

type routingconflicts struct {
	EnvironmentGroup      string                `json:"environmentGroup,omitempty"`
	ConflictingDeployment conflictingdeployment `json:"conflictingDeployment,omitempty"`
	Description           string                `json:"description,omitempty"`
}

type conflictingdeployment struct {
	Basepath    string `json:"basePath,omitempty"`
	Environment string `json:"environment,omitempty"`
	ApiProxy    string `json:"apiProxy,omitempty"`
	Revision    string `json:"revision,omitempty"`
}

// CreateProxy
func CreateProxy(name string, proxy string) (respBody []byte, err error) {
	if proxy != "" {
		err = apiclient.ImportBundle("apis", name, proxy)
		return respBody, err
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis")
	proxyName := "{\"name\":\"" + name + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), proxyName)
	return respBody, err
}

// DeleteProxy
func DeleteProxy(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// DeleteProxyRevision
func DeleteProxyRevision(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name, "revisions", strconv.Itoa(revision))
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// DeployProxy
func DeployProxy(name string, revision int, overrides bool, sequencedRollout bool, safeDeploy bool, serviceAccountName string) (respBody []byte, err error) {

	if safeDeploy {
		var safeResp []byte
		d := deploychangereport{}
		apiclient.SetClientPrintHttpResponse(false)
		if safeResp, err = GenerateDeployChangeReport(name, revision, overrides); err != nil {
			return nil, err
		}
		apiclient.SetClientPrintHttpResponse(apiclient.GetCmdPrintHttpResponseSetting())
		if err = json.Unmarshal(safeResp, &d); err != nil {
			return nil, err
		}
		if len(d.ValidationErrors.Violations) > 0 || len(d.RoutingConflicts) > 0 {
			return nil, fmt.Errorf("Deployment conflicts detected. Here are the conflicts: %s", string(safeResp))
		}
	}

	u, _ := url.Parse(apiclient.BaseURL)

	q := u.Query()
	q.Set("sequencedRollout", strconv.FormatBool(sequencedRollout))

	if overrides || serviceAccountName != "" {
		if overrides {
			q.Set("override", "true")
		}
		if serviceAccountName != "" {
			q.Set("serviceAccount", serviceAccountName)
		}
	}

	u.RawQuery = q.Encode()

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"apis", name, "revisions", strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

// FetchProxy
func FetchProxy(name string, revision int) (err error) {
	return apiclient.FetchBundle("apis", "", name, strconv.Itoa(revision), true)
}

// GetProxy
func GetProxy(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if revision != -1 {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name, "revisions", strconv.Itoa(revision))
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name)
	}
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetHighestProxyRevision
func GetHighestProxyRevision(name string) (version int, err error) {
	apiclient.SetClientPrintHttpResponse(false)
	defer apiclient.SetClientPrintHttpResponse(apiclient.GetCmdPrintHttpResponseSetting())

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name)
	respBody, err := apiclient.HttpClient(u.String())
	if err != nil {
		return -1, err
	}

	proxyRevisions := proxy{}
	if err = json.Unmarshal(respBody, &proxyRevisions); err != nil {
		return -1, err
	}
	version, err = strconv.Atoi(maxRevision(proxyRevisions.Revision))
	if err != nil {
		return -1, nil
	}
	return version, nil
}

// GenerateDeployChangeReport
func GenerateDeployChangeReport(name string, revision int, overrides bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}

	if overrides {
		q := u.Query()
		if overrides {
			q.Set("override", "true")
		}
		u.RawQuery = q.Encode()
	}

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "apis", name, "revisions",
		strconv.Itoa(revision), "deployments:generateDeployChangeReport")
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

// GenerateUndeployChangeReport
func GenerateUndeployChangeReport(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "apis", name, "revisions",
		strconv.Itoa(revision), "deployments:generateUndeployChangeReport")
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

// ListProxies
func ListProxies(includeRevisions bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if includeRevisions {
		q := u.Query()
		q.Set("includeRevisions", strconv.FormatBool(includeRevisions))
		u.RawQuery = q.Encode()
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListEnvDeployments
func ListEnvDeployments() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "deployments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListProxyDeployments
func ListProxyDeployments(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name, "deployments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListProxyRevisionDeployments
func ListProxyRevisionDeployments(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "apis", name, "revisions",
		strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Undeployproxy
func UndeployProxy(name string, revision int, safeUndeploy bool) (respBody []byte, err error) {
	if safeUndeploy {
		var safeResp []byte
		apiclient.SetClientPrintHttpResponse(false)
		if safeResp, err = GenerateUndeployChangeReport(name, revision); err != nil {
			return nil, err
		}
		apiclient.SetClientPrintHttpResponse(apiclient.GetCmdPrintHttpResponseSetting())
		d := deploychangereport{}
		if err = json.Unmarshal(safeResp, &d); err != nil {
			return nil, err
		}
		if len(d.ValidationErrors.Violations) > 0 || len(d.RoutingConflicts) > 0 {
			return nil, fmt.Errorf("Undeployment conflicts detected. Here are the conflicts: %s", string(safeResp))
		}
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"apis", name, "revisions", strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Update
func Update(name string, labels map[string]string) (respBody []byte, err error) {
	if len(labels) != 0 {
		u, _ := url.Parse(apiclient.BaseURL)
		q := u.Query()
		q.Set("updateMask", "labels")
		u.RawQuery = q.Encode()

		labelsArr := []string{}
		for key, value := range labels {
			labelsArr = append(labelsArr, "\""+key+"\":\""+value+"\"")
		}
		payload := "{\"labels\":{" + strings.Join(labelsArr, ",") + "}}"

		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name)
		respBody, err = apiclient.HttpClient(u.String(), payload, "PATCH")
	}

	return respBody, err
}

// CleanProxy
func CleanProxy(name string, reportOnly bool, keepList []string) (err error) {
	type deploymentsDef struct {
		Environment    string `json:"environment,omitempty"`
		APIProxy       string `json:"apiProxy,omitempty"`
		Revision       string `json:"revision,omitempty"`
		DeloyStartTime string `json:"deployStartTime,omitempty"`
		BasePath       string `json:"basePath,omitempty"`
	}

	type proxyDeploymentsDef struct {
		Deployments []deploymentsDef `json:"deployments,omitempty"`
	}

	type metaDataDef struct {
		CreatedAt      string `json:"createdAt,omitempty"`
		LastModifiedAt string `json:"lastModifiedAt,omitempty"`
		SubType        string `json:"subType,omitempty"`
	}

	type proxyRevisionsDef struct {
		MetaData metaDataDef `json:"metaData,omitempty"`
		Name     string      `json:"name,omitempty"`
		Revision []string    `json:"revision,omitempty"`
	}

	proxyDeployments := proxyDeploymentsDef{}
	proxyRevisions := proxyRevisionsDef{}

	reportRevisions := make(map[string]bool)
	deployedRevisions := make(map[string]bool)

	var proxyDeploymentsBytes, proxyRevisionsBytes []byte
	var revision int

	// disable printing
	apiclient.SetClientPrintHttpResponse(false)
	defer apiclient.SetClientPrintHttpResponse(apiclient.GetCmdPrintHttpResponseSetting())

	// step 1. get a list of revisions that are deployed.
	if proxyDeploymentsBytes, err = ListProxyDeployments(name); err != nil {
		return err
	}

	if err = json.Unmarshal(proxyDeploymentsBytes, &proxyDeployments); err != nil {
		return err
	}

	if len(proxyDeployments.Deployments) == 0 {
		return fmt.Errorf("no proxy deployments found")
	}

	for _, proxyDeployment := range proxyDeployments.Deployments {
		if !deployedRevisions[proxyDeployment.Revision] {
			deployedRevisions[proxyDeployment.Revision] = true
		}
	}

	clilog.Info.Println("Revisions [" + getRevisions(deployedRevisions) + "] deployed for API Proxy " + name)

	// step 2. get all the revisions for the proxy
	if proxyRevisionsBytes, err = GetProxy(name, -1); err != nil {
		return err
	}

	if err = json.Unmarshal(proxyRevisionsBytes, &proxyRevisions); err != nil {
		return err
	}

	for _, proxyRevision := range proxyRevisions.Revision {

		// the user passed a list of revisions, check if the revision should be preserved
		deleteRevision := keepRevision(proxyRevision, keepList)

		if !isRevisionDeployed(deployedRevisions, proxyRevision) {
			// step 3. clean up proxy revisions that are not deployed
			if reportOnly {
				if !reportRevisions[proxyRevision] && deleteRevision {
					reportRevisions[proxyRevision] = true
				}
			} else {
				if deleteRevision {
					if revision, err = strconv.Atoi(proxyRevision); err != nil {
						return err
					}
					clilog.Info.Println("Deleting revision: " + proxyRevision)
					if _, err = DeleteProxyRevision(name, revision); err != nil {
						return err
					}
				}
			}
		}
	}

	if reportOnly && len(reportRevisions) > 0 {
		clilog.Info.Println("[REPORT]: API Proxy '" + name + "' revisions: " + getRevisions(reportRevisions) + " can be cleaned")
	}

	return nil
}

// ExportProxies
func ExportProxies(conn int, folder string, allRevisions bool) (err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	q := u.Query()
	q.Set("includeRevisions", "true")
	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis")

	// don't print to sysout
	apiclient.SetClientPrintHttpResponse(false)
	defer apiclient.SetClientPrintHttpResponse(apiclient.GetCmdPrintHttpResponseSetting())

	respBody, err := apiclient.HttpClient(u.String())
	if err != nil {
		return err
	}

	var prxs proxies
	if err = json.Unmarshal(respBody, &prxs); err != nil {
		return err
	}

	clilog.Debug.Printf("Found %d API Proxies in the org\n", len(prxs.Proxies))
	clilog.Debug.Printf("Exporting bundles with parallel %d connections\n", conn)

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
		go exportAPIProxies(&fanOutWg, jobChan, folder, errChan)
	}

	for _, proxy := range prxs.Proxies {
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

func exportAPIProxies(wg *sync.WaitGroup, jobs <-chan revision, folder string, errs chan<- error) {
	defer wg.Done()
	for {
		job, ok := <-jobs
		if !ok {
			return
		}
		err := apiclient.FetchBundle("apis", folder, job.name, job.rev, false)
		if err != nil {
			errs <- err
		}
	}
}

// ImportProxies
func ImportProxies(conn int, folder string) error {
	var bundles []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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

	if len(bundles) == 0 {
		clilog.Warning.Println("No proxy bundle zip files were found in the folder")
		return nil
	}

	clilog.Debug.Printf("Found %d proxy bundles in the folder\n", len(bundles))
	clilog.Debug.Printf("Importing proxies with %d parallel connections\n", conn)

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
		go importAPIProxies(&fanOutWg, jobChan, errChan)
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

func importAPIProxies(wg *sync.WaitGroup, jobs <-chan string, errs chan<- error) {
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
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis")

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

func keepRevision(revision string, keepList []string) bool {
	// if there no keeplist, return true
	if len(keepList) == 0 {
		return true
	}

	for _, r := range keepList {
		if r == revision {
			return false
		}
	}
	return true
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

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
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/clilog"
)

type sharedflows struct {
	Proxies []sharedflow `json:"sharedFlows,omitempty"`
}

type sharedflow struct {
	Name     string   `json:"name,omitempty"`
	Revision []string `json:"revision,omitempty"`
}

//Create
func Create(name string, proxy string) (respBody []byte, err error) {
	if proxy != "" {
		err = apiclient.ImportBundle("sharedflows", name, proxy)
		return respBody, err
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows")
	proxyName := "{\"name\":\"" + name + "\"}"
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), proxyName)
	return respBody, err
}

//Get
func Get(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if revision != -1 {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name, "revisions", strconv.Itoa(revision))
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name)
	}
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//Delete
func Delete(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if revision != -1 {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name, "revisions", strconv.Itoa(revision))
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name)
	}
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//List
func List(includeRevisions bool) (respBody []byte, err error) {

	u, _ := url.Parse(apiclient.BaseURL)
	if includeRevisions {
		q := u.Query()
		q.Set("includeRevisions", strconv.FormatBool(includeRevisions))
		u.RawQuery = q.Encode()
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ListEnvDeployments
func ListEnvDeployments() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}

	q := u.Query()
	q.Set("sharedFlows", "true")
	u.RawQuery = q.Encode()

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ListDeployments
func ListDeployments(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name, "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ListRevisionDeployments
func ListRevisionDeployments(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "sharedflows", name, "revisions",
		strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//Deploy
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
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "")
	return respBody, err
}

//Clean
func Clean(name string, reportOnly bool) (err error) {

	type deploymentsDef struct {
		Environment    string `json:"environment,omitempty"`
		ApiProxy       string `json:"apiProxy,omitempty"`
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

	//disable printing
	apiclient.SetPrintOutput(false)

	//step 1. get a list of revisions that are deployed.
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

	fmt.Println("Revisions [" + getRevisions(deployedRevisions) + "] deployed for Sharedflow " + name)

	//step 2. get all the revisions for the sf
	if sfRevisionsBytes, err = Get(name, -1); err != nil {
		return err
	}

	if err = json.Unmarshal(sfRevisionsBytes, &sfRevisions); err != nil {
		return err
	}

	//enable printing
	apiclient.SetPrintOutput(true)

	for _, sfRevision := range sfRevisions.Revision {
		if !isRevisionDeployed(deployedRevisions, sfRevision) {
			//step 3. clean up proxy revisions that are not deployed
			if reportOnly {
				if !reportRevisions[sfRevision] {
					reportRevisions[sfRevision] = true
				}
			} else {
				if revision, err = strconv.Atoi(sfRevision); err != nil {
					return err
				}
				fmt.Println("Deleting revision: " + sfRevision)
				if _, err = Delete(name, revision); err != nil {
					return err
				}
			}
		}
	}

	if reportOnly && len(reportRevisions) > 0 {
		fmt.Println("[REPORT]: Sharedflow '" + name + "' revisions: " + getRevisions(reportRevisions) + " can be cleaned")
	}

	return nil
}

//Undeploy
func Undeploy(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"sharedflows", name, "revisions", strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//Fetch
func Fetch(name string, revision int) (err error) {
	return apiclient.FetchBundle("sharedflows", "", name, strconv.Itoa(revision), true)
}

//Export
func Export(conn int, folder string, allRevisions bool) (err error) {
	//parent workgroup
	var pwg sync.WaitGroup
	const entityType = "sharedflows"

	u, _ := url.Parse(apiclient.BaseURL)
	q := u.Query()
	q.Set("includeRevisions", "true")
	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), entityType)

	//don't print to sysout
	respBody, err := apiclient.HttpClient(false, u.String())
	if err != nil {
		return err
	}

	var entities = sharedflows{}
	err = json.Unmarshal(respBody, &entities)
	if err != nil {
		return err
	}

	numEntities := len(entities.Proxies)
	clilog.Info.Printf("Found %d API Proxies in the org\n", numEntities)
	clilog.Info.Printf("Exporting bundles with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than products
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Exporting batch %d of proxies\n", (i + 1))
		go batchExport(entities.Proxies[start:end], entityType, folder, allRevisions, &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Exporting remaining %d proxies\n", remaining)
		go batchExport(entities.Proxies[start:numEntities], entityType, folder, allRevisions, &pwg)
		pwg.Wait()
	}

	return nil
}

//Import
func Import(conn int, folder string) error {
	var pwg sync.WaitGroup
	var entities []string

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
		entities = append(entities, path)
		return nil
	})

	if err != nil {
		return err
	}

	numEntities := len(entities)
	clilog.Info.Printf("Found %d proxy bundles in the folder\n", numEntities)
	clilog.Info.Printf("Create proxies with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Creating batch %d of bundles\n", (i + 1))
		go batchImport(entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Creating remaining %d bundles\n", remaining)
		go batchImport(entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

func batchExport(entities []sharedflow, entityType string, folder string, allRevisions bool, pwg *sync.WaitGroup) {
	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup
	//proxy revision wait group
	var prwg sync.WaitGroup
	//number of revisions to download concurrently
	conn := 2

	numSharedFlows := len(entities)
	bwg.Add(numSharedFlows)

	for _, entity := range entities {
		if !allRevisions {
			//download only the last revision
			lastRevision := maxRevision(entity.Revision)
			clilog.Info.Printf("Downloading revision %s of sharedflow %s\n", lastRevision, entity.Name)
			go apiclient.FetchAsyncBundle(entityType, folder, entity.Name, lastRevision, allRevisions, &bwg)
		} else {
			//download all revisions
			if len(entity.Revision) == 1 {
				lastRevision := maxRevision(entity.Revision)
				clilog.Info.Printf("Downloading revision %s of sharedflow %s\n", lastRevision, entity.Name)
				go apiclient.FetchAsyncBundle(entityType, folder, entity.Name, lastRevision, allRevisions, &bwg)
			} else {

				numRevisions := len(entity.Revision)
				numOfLoops, remaining := numRevisions/conn, numRevisions%conn

				start := 0

				for i, end := 0, 0; i < numOfLoops; i++ {
					prwg.Add(1)
					end = (i * conn) + conn
					clilog.Info.Printf("Exporting batch %d of sharedflow revisions\n", (i + 1))
					go batchExportRevisions(entityType, folder, entity.Name, entity.Revision[start:end], &prwg)
					start = end
					prwg.Wait()
				}

				if remaining > 0 {
					prwg.Add(1)
					clilog.Info.Printf("Exporting remaining %d sharedflow revisions\n", remaining)
					go batchExportRevisions(entityType, folder, entity.Name, entity.Revision[start:numRevisions], &prwg)
					prwg.Wait()
				}
				bwg.Done()
			}
		}
	}
	bwg.Wait()
}

//batchExportRevisions
func batchExportRevisions(entityType string, folder string, name string, revisions []string, prwg *sync.WaitGroup) {
	defer prwg.Done()
	//batch workgroup
	var brwg sync.WaitGroup

	brwg.Add(len(revisions))

	for _, revision := range revisions {
		clilog.Info.Printf("Exporting sharedflow %s revision %s\n", name, revision)
		go apiclient.FetchAsyncBundle(entityType, folder, name, revision, true, &brwg)
	}
	brwg.Wait()
}

//batch creates a batch of sharedflow to import
func batchImport(entities []string, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		//sharedflow name is empty; same as filename
		go apiclient.ImportBundleAsync("sharedflows", "", entity, &bwg)
	}
	bwg.Wait()
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

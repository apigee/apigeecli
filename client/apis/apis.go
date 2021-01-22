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

type proxies struct {
	Proxies []proxy `json:"proxies,omitempty"`
}

type proxy struct {
	Name     string   `json:"name,omitempty"`
	Revision []string `json:"revision,omitempty"`
}

//CreateProxy
func CreateProxy(name string, proxy string) (respBody []byte, err error) {
	if proxy != "" {
		err = apiclient.ImportBundle("apis", name, proxy)
		return respBody, err
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis")
	proxyName := "{\"name\":\"" + name + "\"}"
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), proxyName)
	return respBody, err
}

//DeleteProxy
func DeleteProxy(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if revision != -1 {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name, "revisions", strconv.Itoa(revision))
	}
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//DeployProxy
func DeployProxy(name string, revision int, overrides bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if overrides {
		q := u.Query()
		q.Set("override", "true")
		u.RawQuery = q.Encode()
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"apis", name, "revisions", strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "")
	return respBody, err
}

//FetchProxy
func FetchProxy(name string, revision int) (err error) {
	return apiclient.FetchBundle("apis", "", name, strconv.Itoa(revision))
}

//GetProxy
func GetProxy(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if revision != -1 {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name, "revisions", strconv.Itoa(revision))
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name)
	}
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ListProxies
func ListProxies(includeRevisions bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if includeRevisions {
		q := u.Query()
		q.Set("includeRevisions", strconv.FormatBool(includeRevisions))
		u.RawQuery = q.Encode()
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ListEnvDeployments
func ListEnvDeployments() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ListProxyDeployments
func ListProxyDeployments(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name, "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ListProxyRevisionDeployments
func ListProxyRevisionDeployments(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() == "" {
		return respBody, fmt.Errorf("environment name missing")
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "apis", name, "revisions",
		strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//Undeployproxy
func UndeployProxy(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"apis", name, "revisions", strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//CleanProxy
func CleanProxy(name string, reportOnly bool) (err error) {

	type deploymentsDef struct {
		Environment    string `json:"environment,omitempty"`
		ApiProxy       string `json:"apiProxy,omitempty"`
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

	//disable printing
	apiclient.SetPrintOutput(false)

	//step 1. get a list of revisions that are deployed.
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

	fmt.Println("Revisions [" + getRevisions(deployedRevisions) + "] deployed for API Proxy " + name)

	//step 2. get all the revisions for the proxy
	if proxyRevisionsBytes, err = GetProxy(name, -1); err != nil {
		return err
	}

	if err = json.Unmarshal(proxyRevisionsBytes, &proxyRevisions); err != nil {
		return err
	}

	//enable printing
	apiclient.SetPrintOutput(true)

	for _, proxyRevision := range proxyRevisions.Revision {
		if !isRevisionDeployed(deployedRevisions, proxyRevision) {
			//step 3. clean up proxy revisions that are not deployed
			if reportOnly {
				if !reportRevisions[proxyRevision] {
					reportRevisions[proxyRevision] = true
				}
			} else {
				if revision, err = strconv.Atoi(proxyRevision); err != nil {
					return err
				}
				fmt.Println("Deleting revision: " + proxyRevision)
				if _, err = DeleteProxy(name, revision); err != nil {
					return err
				}
			}
		}
	}

	if reportOnly && len(reportRevisions) > 0 {
		fmt.Println("[REPORT]: API Proxy '" + name + "' revisions: " + getRevisions(reportRevisions) + " can be cleaned")
	}

	return nil
}

//ExportProxies
func ExportProxies(conn int, folder string) (err error) {
	//parent workgroup
	var pwg sync.WaitGroup
	const entityType = "apis"

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

	var entities = proxies{}
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
		go batchExport(entities.Proxies[start:end], entityType, folder, &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Exporting remaining %d proxies\n", remaining)
		go batchExport(entities.Proxies[start:numEntities], entityType, folder, &pwg)
		pwg.Wait()
	}

	return nil
}

//ImportProxies
func ImportProxies(conn int, folder string) error {

	var pwg sync.WaitGroup
	var entities []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
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

func batchExport(entities []proxy, entityType string, folder string, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		//download only the last revision
		lastRevision := len(entity.Revision)
		go apiclient.FetchAsyncBundle(entityType, folder, entity.Name, entity.Revision[lastRevision-1], &bwg)
	}
	bwg.Wait()
}

//batch creates a batch of proxies to import
func batchImport(entities []string, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		//api proxy name is empty; same as filename
		go apiclient.ImportBundleAsync("apis", "", entity, &bwg)
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

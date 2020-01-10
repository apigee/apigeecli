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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

func List(includeRevisions bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() != "" {
		q := u.Query()
		q.Set("sharedFlows", "true")
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "deployments")
	} else {
		if includeRevisions {
			q := u.Query()
			q.Set("includeRevisions", strconv.FormatBool(includeRevisions))
			u.RawQuery = q.Encode()
		}
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sharedflows")
	}
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

func Deploy(name string, revision int, overrides bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if overrides {
		q := u.Query()
		q.Set("override", "true")
		u.RawQuery = q.Encode()
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"sharedflows", name, "revisions", strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "")
	return respBody, err
}

func Undeploy(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"sharedflows", name, "revisions", strconv.Itoa(revision), "deployments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

func Fetch(name string, revision int) (err error) {
	return apiclient.FetchBundle("sharedflows", name, strconv.Itoa(revision))
}

func Export(conn int) (err error) {
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
		go batchExport(entities.Proxies[start:end], entityType, &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Exporting remaining %d proxies\n", remaining)
		go batchExport(entities.Proxies[start:numEntities], entityType, &pwg)
		pwg.Wait()
	}

	return nil
}

func Import(conn int, folder string) error {
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

func batchExport(entities []sharedflow, entityType string, pwg *sync.WaitGroup) {
	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		//download only the last revision
		lastRevision := len(entity.Revision)
		go apiclient.FetchAsyncBundle(entityType, entity.Name, entity.Revision[lastRevision-1], &bwg)
	}
	bwg.Wait()
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

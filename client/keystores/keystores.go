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

package keystores

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/clilog"
)

// Create
func Create(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores")
	payload := "{\"name\":\"" + name + "\"}"
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

// List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

// Import
func Import(conn int, filePath string) (err error) {
	var pwg sync.WaitGroup
	const entityType = "keystores"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), entityType)

	entities, err := readKeystoresFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	clilog.Info.Printf("Found %d keystores in the file\n", numEntities)
	clilog.Info.Printf("Create keystores with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Creating batch %d of keystores\n", (i + 1))
		go batchImport(u.String(), entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Creating remaining %d keystores\n", remaining)
		go batchImport(u.String(), entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

// batchImport creates a batch of keystores to create
func batchImport(url string, entities []string, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncKeystore(url, entity, &bwg)
	}
	bwg.Wait()
}

func createAsyncKeystore(url string, entity string, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := Create(entity)
	if err != nil {
		clilog.Error.Println(err)
		return
	}
	clilog.Info.Printf("Completed entity: %s", entity)
}

func readKeystoresFile(filePath string) ([]string, error) {
	keystoresList := []string{}

	jsonFile, err := os.Open(filePath)

	if err != nil {
		return keystoresList, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		return keystoresList, err
	}

	err = json.Unmarshal(byteValue, &keystoresList)

	if err != nil {
		return keystoresList, err
	}

	return keystoresList, nil
}

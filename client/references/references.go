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

package references

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/clilog"
)

type ref struct {
	Name         string `json:"name,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	Refers       string `json:"refers,omitempty"`
}

// Create references
func Create(name string, description string, resourceType string, refers string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	reference := []string{}

	reference = append(reference, "\"name\":\""+name+"\"")

	if description != "" {
		reference = append(reference, "\"description\":\""+description+"\"")
	}

	reference = append(reference, "\"resourceType\":\""+resourceType+"\"")
	reference = append(reference, "\"refers\":\""+refers+"\"")

	payload := "{" + strings.Join(reference, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

// Get a reference
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

// Delete a reference
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

// List references
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

// Update references
func Update(name string, description string, resourceType string, refers string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	reference := []string{}

	reference = append(reference, "\"name\":\""+name+"\"")

	if description != "" {
		reference = append(reference, "\"description\":\""+description+"\"")
	}

	if resourceType != "" {
		reference = append(reference, "\"resourceType\":\""+resourceType+"\"")
	}

	if refers != "" {
		reference = append(reference, "\"refers\":\""+refers+"\"")
	}

	payload := "{" + strings.Join(reference, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload, "PUT")
	return respBody, err
}

// Export
func Export(conn int) (payload [][]byte, err error) {
	//parent workgroup
	var pwg sync.WaitGroup
	var mu sync.Mutex
	const entityType = "references"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), entityType)
	//don't print to sysout
	respBody, err := apiclient.HttpClient(false, u.String())
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	var refList []string
	err = json.Unmarshal(respBody, &refList)
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	numEntities := len(refList)
	clilog.Info.Printf("Found %d references in the org\n", numEntities)
	clilog.Info.Printf("Exporting references with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than references
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Exporting batch %d of references\n", (i + 1))
		go batchExport(refList[start:end], entityType, &pwg, &mu)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Exporting remaining %d references\n", remaining)
		go batchExport(refList[start:numEntities], entityType, &pwg, &mu)
		pwg.Wait()
	}

	payload = make([][]byte, len(apiclient.GetEntityPayloadList()))
	copy(payload, apiclient.GetEntityPayloadList())
	apiclient.ClearEntityPayloadList()
	return payload, nil
}

// batch created a batch of references to query
func batchExport(entities []string, entityType string, pwg *sync.WaitGroup, mu *sync.Mutex) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		u, _ := url.Parse(apiclient.BaseURL)
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), entityType, entity)
		go apiclient.GetAsyncEntity(u.String(), &bwg, mu)
	}
	bwg.Wait()
}

// Import
func Import(conn int, filePath string) (err error) {
	var pwg sync.WaitGroup
	const entityType = "references"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), entityType)

	entities, err := readReferencesFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	clilog.Info.Printf("Found %d references in the file\n", numEntities)
	clilog.Info.Printf("Create references with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Creating batch %d of references\n", (i + 1))
		go batchImport(u.String(), entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Creating remaining %d references\n", remaining)
		go batchImport(u.String(), entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

// batchImport creates a batch of references to create
func batchImport(url string, entities []ref, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncReference(url, entity, &bwg)
	}
	bwg.Wait()
}

func createAsyncReference(url string, entity ref, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := json.Marshal(entity)
	if err != nil {
		clilog.Error.Println(err)
		return
	}
	_, err = apiclient.HttpClient(apiclient.GetPrintOutput(), url, string(out))
	if err != nil {
		clilog.Error.Println(err)
		return
	}
	clilog.Info.Printf("Completed entity: %s", entity.Name)
}

func readReferencesFile(filePath string) ([]ref, error) {
	refList := []ref{}

	jsonFile, err := os.Open(filePath)

	if err != nil {
		return refList, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		return refList, err
	}

	err = json.Unmarshal(byteValue, &refList)

	if err != nil {
		return refList, err
	}

	return refList, nil
}

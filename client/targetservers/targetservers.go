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

package targetservers

import (
	"encoding/json"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/clilog"
)

//Create
func Create(name string, description string, host string, port int, enable bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	targetserver := []string{}

	targetserver = append(targetserver, "\"name\":\""+name+"\"")

	if description != "" {
		targetserver = append(targetserver, "\"description\":\""+description+"\"")
	}

	targetserver = append(targetserver, "\"host\":\""+host+"\"")
	targetserver = append(targetserver, "\"port\":"+strconv.Itoa(port))

	if enable {
		targetserver = append(targetserver, "\"isEnabled\":"+strconv.FormatBool(enable))
	}

	payload := "{" + strings.Join(targetserver, ",") + "}"
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//Export
func Export(conn int) (payload [][]byte, err error) {

	//parent workgroup
	var pwg sync.WaitGroup
	var mu sync.Mutex
	const entityType = "targetservers"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), entityType)
	//don't print to sysout
	respBody, err := apiclient.HttpClient(false, u.String())
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	var targetservers []string
	err = json.Unmarshal(respBody, &targetservers)
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	numEntities := len(targetservers)
	clilog.Info.Printf("Found %d targetservers in the org\n", numEntities)
	clilog.Info.Printf("Exporting targetservers with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than targetservers
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Exporting batch %d of targetservers\n", (i + 1))
		go batchExport(targetservers[start:end], entityType, &pwg, &mu)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Exporting remaining %d targetservers\n", remaining)
		go batchExport(targetservers[start:numEntities], entityType, &pwg, &mu)
		pwg.Wait()
	}

	return apiclient.GetEntityPayloadList(), nil
}

//batch created a batch of targetservers to query
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

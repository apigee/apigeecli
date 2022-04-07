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
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/clilog"
)

type targetserver struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Host        string   `json:"host,omitempty"`
	Port        int      `json:"port,omitempty"`
	IsEnabled   bool     `json:"isEnabled,omitempty"`
	Protocol    string   `json:"protocol,omitempty"`
	SslInfo     *sslInfo `json:"sSLInfo,omitempty"`
}

type sslInfo struct {
	Enabled                bool        `json:"enabled,omitempty"`
	ClientAuthEnabled      bool        `json:"clientAuthEnabled,omitempty"`
	Keystore               string      `json:"keyStore,omitempty"`
	Keyalias               string      `json:"keyAlias,omitempty"`
	Truststore             string      `json:"trustStore,omitempty"`
	IgnoreValidationErrors bool        `json:"ignoreValidationErrors,omitempty"`
	Protocols              []string    `json:"protocols,omitempty"`
	Ciphers                []string    `json:"ciphers,omitempty"`
	CommonName             *commonName `json:"commonName,omitempty"`
}

type commonName struct {
	Value         string `json:"value,omitempty"`
	WildcardMatch bool   `json:"wildcardMatch,omitempty"`
}

//Create
func Create(name string, description string, host string, port int, enable bool, grpc bool, keyStore string, keyAlias string, sslinfo bool, tlsenabled bool, clientAuthEnabled bool, ignoreValidationErrors bool) (respBody []byte, err error) {
	targetsvr := targetserver{}
	targetsvr.Name = name

	return createOrUpdate("create", targetsvr, name, description, host, port, enable, grpc, keyStore, keyAlias, sslinfo, tlsenabled, clientAuthEnabled, ignoreValidationErrors)
}

//Update
func Update(name string, description string, host string, port int, enable bool, grpc bool, keyStore string, keyAlias string, sslinfo bool, tlsenabled bool, clientAuthEnabled bool, ignoreValidationErrors bool) (respBody []byte, err error) {

	var targetRespBody []byte
	var targetsvr = targetserver{}

	apiclient.SetPrintOutput(false)
	if targetRespBody, err = Get(name); err != nil {
		return nil, err
	}
	apiclient.SetPrintOutput(true)

	if err = json.Unmarshal(targetRespBody, &targetsvr); err != nil {
		return nil, err
	}

	return createOrUpdate("update", targetsvr, name, description, host, port, enable, grpc, keyStore, keyAlias, sslinfo, tlsenabled, clientAuthEnabled, ignoreValidationErrors)
}

func createOrUpdate(action string, targetsvr targetserver, name string, description string, host string, port int, enable bool, grpc bool, keyStore string, keyAlias string, sslinfo bool, tlsenabled bool, clientAuthEnabled bool, ignoreValidationErrors bool) (respBody []byte, err error) {

	var reqBody []byte
	sslInfoObj := new(sslInfo)

	if description != "" {
		targetsvr.Description = description
	}

	if host != "" {
		targetsvr.Host = host
	}

	if port != -1 {
		targetsvr.Port = port
	}

	if grpc {
		targetsvr.Protocol = "GRPC"
	}

	if enable {
		targetsvr.IsEnabled = enable
	}

	if tlsenabled {
		sslInfoObj.Enabled = tlsenabled
	}

	if clientAuthEnabled {
		sslInfoObj.ClientAuthEnabled = clientAuthEnabled
	}

	if ignoreValidationErrors {
		sslInfoObj.IgnoreValidationErrors = ignoreValidationErrors
	}

	if keyAlias != "" {
		sslInfoObj.Keyalias = keyAlias
	}

	if keyStore != "" {
		sslInfoObj.Keystore = keyStore
	}

	if sslinfo {
		targetsvr.SslInfo = sslInfoObj
	}

	if reqBody, err = json.Marshal(targetsvr); err != nil {
		return nil, err
	}

	u, _ := url.Parse(apiclient.BaseURL)

	if action == "create" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers")
		respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(reqBody))
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers", name)
		respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(reqBody), "PUT")
	}

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

	payload = make([][]byte, len(apiclient.GetEntityPayloadList()))
	copy(payload, apiclient.GetEntityPayloadList())
	apiclient.ClearEntityPayloadList()
	return payload, nil
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

//Import
func Import(conn int, filePath string) (err error) {
	var pwg sync.WaitGroup
	const entityType = "targetservers"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), entityType)

	entities, err := readTargetServersFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	clilog.Info.Printf("Found %d target servers in the file\n", numEntities)
	clilog.Info.Printf("Create target servers with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Creating batch %d of target servers\n", (i + 1))
		go batchImport(u.String(), entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Creating remaining %d target servers\n", remaining)
		go batchImport(u.String(), entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

//batch creates a batch of target servers to create
func batchImport(url string, entities []targetserver, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncTargetServer(url, entity, &bwg)
	}
	bwg.Wait()
}

func createAsyncTargetServer(url string, entity targetserver, wg *sync.WaitGroup) {
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

func readTargetServersFile(filePath string) ([]targetserver, error) {

	targetservers := []targetserver{}

	jsonFile, err := os.Open(filePath)

	if err != nil {
		return targetservers, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return targetservers, err
	}

	err = json.Unmarshal(byteValue, &targetservers)

	if err != nil {
		return targetservers, err
	}

	return targetservers, nil
}

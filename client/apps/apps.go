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

package apps

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/clilog"
	"github.com/thedevsaddam/gojsonq"
)

type apps struct {
	Apps []app `json:"app,omitempty"`
}

type app struct {
	AppID string `json:"appId,omitempty"`
}

type application struct {
	Name        string        `json:"name,omitempty"`
	Status      string        `json:"status,omitempty"`
	Credentials *[]credential `json:"credentials,omitempty"`
	DeveloperID *string       `json:"developerId,omitempty"`
	DisplayName string        `json:"displayName,omitempty"`
	Attributes  []attribute   `json:"attributes,omitempty"`
	CallbackURL string        `json:"callbackUrl,omitempty"`
	Scopes      []string      `json:"scopes,omitempty"`
}

type credential struct {
	APIProducts    []apiProduct `json:"apiProducts,omitempty"`
	ConsumerKey    string       `json:"consumerKey,omitempty"`
	ConsumerSecret string       `json:"consumerSecret,omitempty"`
	ExpiresAt      int          `json:"expiresAt,omitempty"`
	Status         string       `json:"status,omitempty"`
	Scopes         []string     `json:"scopes,omitempty"`
}

type apiProduct struct {
	Name string `json:"apiproduct,omitempty"`
}

type importCredential struct {
	APIProducts    []string `json:"apiProducts,omitempty"`
	ConsumerKey    string   `json:"consumerKey,omitempty"`
	ConsumerSecret string   `json:"consumerSecret,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
}

//attribute to used to hold custom attributes for entities
type attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func Create(name string, email string, expires string, callback string, apiProducts []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	app := []string{}

	app = append(app, "\"name\":\""+name+"\"")

	if len(apiProducts) > 0 {
		app = append(app, "\"apiProducts\":[\""+getArrayStr(apiProducts)+"\"]")
	}

	if callback != "" {
		app = append(app, "\"callbackUrl\":\""+callback+"\"")
	}

	if expires != "" {
		app = append(app, "\"keyExpiresIn\":\""+expires+"\"")
	}

	if len(scopes) > 0 {
		app = append(app, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")
	}

	if len(attrs) != 0 {
		attributes := []string{}
		for key, value := range attrs {
			attributes = append(attributes, "{\"name\":\""+key+"\",\"value\":\""+value+"\"}")
		}
		attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
		app = append(app, attributesStr)
	}

	payload := "{" + strings.Join(app, ",") + "}"
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", email, "apps")
	respBody, err = apiclient.HttpClient(true, u.String(), payload)
	return respBody, err
}

func Delete(name string, developerID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerID, "apps", name)
	respBody, err = apiclient.HttpClient(true, u.String(), "", "DELETE")
	return respBody, err
}

func Get(appID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apps", appID)
	respBody, err = apiclient.HttpClient(true, u.String())
	return respBody, err
}

func SearchApp(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	//search by name is not implemented; use list and return the appropriate app
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apps")
	q := u.Query()
	q.Set("expand", "true")
	q.Set("includeCred", "false")
	u.RawQuery = q.Encode()
	//don't print the list
	respBody, err = apiclient.HttpClient(false, u.String())
	if err != nil {
		return respBody, err
	}
	jq := gojsonq.New().JSONString(string(respBody)).From("app").Where("name", "eq", name)
	out := jq.Get()
	outBytes, err := json.Marshal(out)
	if err != nil {
		return outBytes, err
	}
	return outBytes, nil
}

func List(includeCred bool, expand bool, count int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apps")
	q := u.Query()
	if expand {
		q.Set("expand", "true")
	} else {
		q.Set("expand", "false")
	}
	if expand && includeCred {
		q.Set("includeCred", "true")
	} else if expand && !includeCred {
		q.Set("includeCred", "false")
	}
	if count != -1 {
		q.Set("row", strconv.Itoa(count))
	}
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(true, u.String())
	return respBody, err
}

func GenerateKey(name string, developerID string, apiProducts []string, callback string, expires string, scopes []string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	key := []string{}

	key = append(key, "\"name\":\""+name+"\"")
	key = append(key, "\"apiProducts\":[\""+getArrayStr(apiProducts)+"\"]")

	if callback != "" {
		key = append(key, "\"callbackUrl\":\""+callback+"\"")
	}

	if expires != "" {
		key = append(key, "\"keyExpiresIn\":\""+expires+"\"")
	}

	if len(scopes) > 0 {
		key = append(key, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")
	}

	payload := "{" + strings.Join(key, ",") + "}"
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerID, "apps", name)
	respBody, err = apiclient.HttpClient(true, u.String(), payload)
	return respBody, err
}

func Export(conn int) (payload [][]byte, err error) {
	//parent workgroup
	var pwg sync.WaitGroup
	var mu sync.Mutex
	const entityType = "apps"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), entityType)
	//don't print to sysout
	respBody, err := apiclient.HttpClient(false, u.String())
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	var entities = apps{}
	err = json.Unmarshal(respBody, &entities)
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	numEntities := len(entities.Apps)
	clilog.Info.Printf("Found %d apps in the org\n", numEntities)
	clilog.Info.Printf("Exporting apps with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than apps
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Exporting batch %d of apps\n", (i + 1))
		go batchExport(entities.Apps[start:end], entityType, &pwg, &mu)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Exporting remaining %d apps\n", remaining)
		go batchExport(entities.Apps[start:numEntities], entityType, &pwg, &mu)
		pwg.Wait()
	}

	return apiclient.GetEntityPayloadList(), nil
}

func Import(conn int, filePath string) error {
	var pwg sync.WaitGroup

	entities, err := readAppsFile(filePath)
	if err != nil {
		clilog.Error.Println("error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	clilog.Info.Printf("Found %d apps in the file\n", numEntities)
	clilog.Info.Printf("Create apps with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Creating batch %d of apps\n", (i + 1))
		go batchImport(entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Creating remaining %d apps\n", remaining)
		go batchImport(entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

func readAppsFile(filePath string) ([]application, error) {

	apps := []application{}

	jsonFile, err := os.Open(filePath)

	if err != nil {
		return apps, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return apps, err
	}

	err = json.Unmarshal(byteValue, &apps)

	if err != nil {
		return apps, err
	}

	return apps, nil
}

//batch created a batch of apps to query
func batchExport(entities []app, entityType string, pwg *sync.WaitGroup, mu *sync.Mutex) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		u, _ := url.Parse(apiclient.BaseURL)
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), entityType, entity.AppID)
		go apiclient.GetAsyncEntity(u.String(), &bwg, mu)
	}
	bwg.Wait()
}

//batch creates a batch of apps to create
func batchImport(entities []application, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncApp(entity, &bwg)
	}
	bwg.Wait()
}

func createAsyncApp(app application, wg *sync.WaitGroup) {
	defer wg.Done()

	//importing an app will be a two step process.
	//1. create the app without the credential
	//2. create/import the credential
	u, _ := url.Parse(apiclient.BaseURL)
	//store the developer and the credential
	developerID := *app.DeveloperID
	credentials := *app.Credentials

	//remove the developer id and credentials from the payload
	app.DeveloperID = nil
	app.Credentials = nil

	out, err := json.Marshal(app)
	if err != nil {
		clilog.Error.Println(err)
		return
	}

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerID, "apps")
	_, err = apiclient.HttpClient(true, u.String(), string(out))
	if err != nil {
		clilog.Error.Println(err)
		return
	}
	u, _ = url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerID, "apps", app.Name, "keys", "create")
	for _, credential := range credentials {
		//construct a []string for products
		var products []string
		for _, apiProduct := range credential.APIProducts {
			products = append(products, apiProduct.Name)
		}
		//create a new credential
		importCred := importCredential{}
		importCred.APIProducts = products
		importCred.ConsumerKey = credential.ConsumerKey
		importCred.ConsumerSecret = credential.ConsumerSecret
		importCred.Scopes = credential.Scopes

		impCredJSON, err := json.Marshal(importCred)
		if err != nil {
			clilog.Error.Println(err)
			return
		}
		_, err = apiclient.HttpClient(true, u.String(), string(impCredJSON))
		if err != nil {
			return
		}
		clilog.Warning.Println("NOTE: apiProducts are not associated with the app")
	}
	clilog.Info.Printf("Completed entity: %s", app.Name)
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}

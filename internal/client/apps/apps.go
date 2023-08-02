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
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"internal/apiclient"

	"internal/clilog"

	"internal/client/developers"

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
	ExpiresAt      string       `json:"expiresAt,omitempty"`
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

// attribute to used to hold custom attributes for entities
type attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Action uint8

const (
	CREATE Action = iota
	UPDATE
)

// Create
func Create(name string, email string, expires string, callback string, apiProducts []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	return createOrUpdate(name, email, expires, callback, apiProducts, scopes, attrs, CREATE)
}

// Delete
func Delete(name string, developerID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerID, "apps", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get(appID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apps", appID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Update
func Update(name string, email string, expires string, callback string, apiProducts []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	return createOrUpdate(name, email, expires, callback, apiProducts, scopes, attrs, UPDATE)
}

func createOrUpdate(name string, email string, expires string, callback string, apiProducts []string, scopes []string, attrs map[string]string, action Action) (respBody []byte, err error) {
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

	if action == CREATE {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", email, "apps")
		respBody, err = apiclient.HttpClient(u.String(), payload)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", email, "apps", name)
		respBody, err = apiclient.HttpClient(u.String(), payload, "PUT")
	}

	return respBody, err
}

// Manage
func Manage(appID string, developerEmail string, action string) (respBody []byte, err error) {
	if action != "revoke" && action != "approve" {
		return nil, fmt.Errorf("invalid action. action must be revoke or approve")
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerEmail, "apps", appID)
	q := u.Query()
	q.Set("action", action)
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), "", "POST", "application/octet-stream")
	return respBody, err
}

// SearchApp
func SearchApp(name string) (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	u, _ := url.Parse(apiclient.BaseURL)
	// search by name is not implemented; use list and return the appropriate app
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apps")
	q := u.Query()
	q.Set("expand", "true")
	q.Set("includeCred", "false")
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String())
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

// List
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
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListApps
func ListApps(productName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apps")
	q := u.Query()
	q.Set("apiProduct", productName)
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GenerateKey
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
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Export
func Export(conn int) (payload [][]byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	// parent workgroup
	var pwg sync.WaitGroup
	var mu sync.Mutex
	const entityType = "apps"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), entityType)

	respBody, err := apiclient.HttpClient(u.String())
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	entities := apps{}
	err = json.Unmarshal(respBody, &entities)
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	numEntities := len(entities.Apps)
	clilog.Debug.Printf("Found %d apps in the org\n", numEntities)
	clilog.Debug.Printf("Exporting apps with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	// ensure connections aren't greater than apps
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Debug.Printf("Exporting batch %d of apps\n", (i + 1))
		go batchExport(entities.Apps[start:end], entityType, &pwg, &mu)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Debug.Printf("Exporting remaining %d apps\n", remaining)
		go batchExport(entities.Apps[start:numEntities], entityType, &pwg, &mu)
		pwg.Wait()
	}
	payload = make([][]byte, len(apiclient.GetEntityPayloadList()))
	copy(payload, apiclient.GetEntityPayloadList())
	apiclient.ClearEntityPayloadList()
	return payload, nil
}

// Import
func Import(conn int, filePath string, developersFilePath string) error {
	var pwg sync.WaitGroup

	entities, developerEntities, err := readAppsFile(filePath, developersFilePath)
	if err != nil {
		clilog.Error.Println("error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	clilog.Debug.Printf("Found %d apps in the file\n", numEntities)
	clilog.Debug.Printf("Create apps with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	// ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Debug.Printf("Creating batch %d of apps\n", (i + 1))
		go batchImport(entities[start:end], developerEntities, &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Debug.Printf("Creating remaining %d apps\n", remaining)
		go batchImport(entities[start:numEntities], developerEntities, &pwg)
		pwg.Wait()
	}

	return nil
}

func readAppsFile(filePath string, developersFilePath string) ([]application, developers.Appdevelopers, error) {
	apps := []application{}
	devs := developers.Appdevelopers{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return apps, devs, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return apps, devs, err
	}

	err = json.Unmarshal(byteValue, &apps)
	if err != nil {
		return apps, devs, err
	}

	devs, err = developers.ReadDevelopersFile(developersFilePath)
	if err != nil {
		return apps, devs, err
	}

	return apps, devs, nil
}

// batch created a batch of apps to query
func batchExport(entities []app, entityType string, pwg *sync.WaitGroup, mu *sync.Mutex) {
	defer pwg.Done()
	// batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		u, _ := url.Parse(apiclient.BaseURL)
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), entityType, entity.AppID)
		go apiclient.GetAsyncEntity(u.String(), &bwg, mu)
	}
	bwg.Wait()
}

// batch creates a batch of apps to create
func batchImport(entities []application, developerEntities developers.Appdevelopers, pwg *sync.WaitGroup) {
	defer pwg.Done()
	// batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncApp(entity, developerEntities, &bwg)
	}
	bwg.Wait()
}

func createAsyncApp(app application, developerEntities developers.Appdevelopers, wg *sync.WaitGroup) {
	defer wg.Done()

	// importing an app will be a two step process.
	// 1. create the app without the credential
	// 2. create/import the credential
	u, _ := url.Parse(apiclient.BaseURL)
	if app.DeveloperID == nil {
		clilog.Error.Println("developer id was not found")
		return
	}
	// store the developer and the credential
	developerEmail, developerID, err := getNewDeveloperId(*app.DeveloperID, developerEntities) //*app.DeveloperID
	if err != nil {
		clilog.Error.Println(err)
		return
	}

	credentials := *app.Credentials

	// remove the developer id and credentials from the payload
	app.DeveloperID = nil
	app.Credentials = nil

	out, err := json.Marshal(app)
	if err != nil {
		clilog.Error.Println(err)
		return
	}

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerID, "apps")
	appRespBody, err := apiclient.HttpClient(u.String(), string(out))
	if err != nil {
		clilog.Error.Println(err)
		return
	}

	// get the new appId & keyId
	var newDeveloperApp map[string]interface{}
	err = json.Unmarshal(appRespBody, &newDeveloperApp)
	if err != nil {
		clilog.Error.Println(err)
		return
	}

	// delete the auto-generated key
	newAppCredentials := newDeveloperApp["credentials"].([]interface{})
	temporaryCredential := newAppCredentials[0].(map[string]interface{})

	apiclient.ClientPrintHttpResponse.Set(false)
	_, err = DeleteKey(developerEmail, newDeveloperApp["name"].(string), temporaryCredential["consumerKey"].(string))
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	if err != nil {
		clilog.Error.Println(err)
		return
	}

	createDeveloperAppUrl, _ := url.Parse(apiclient.BaseURL)
	createDeveloperAppUrl.Path = path.Join(createDeveloperAppUrl.Path, apiclient.GetApigeeOrg(), "developers", developerID, "apps", app.Name, "keys")
	for _, credential := range credentials {

		// create a new credential
		importCred := importCredential{}
		importCred.ConsumerKey = credential.ConsumerKey
		importCred.ConsumerSecret = credential.ConsumerSecret

		impCredJSON, err := json.Marshal(importCred)
		if err != nil {
			clilog.Error.Println(err)
			return
		}

		_, err = apiclient.HttpClient(createDeveloperAppUrl.String(), string(impCredJSON))
		if err != nil {
			return
		}

		// update credentials

		// construct a []string for products
		var products []string
		for _, apiProduct := range credential.APIProducts {
			products = append(products, apiProduct.Name)
		}

		if len(products) > 0 {
			// updateDeveloperApp
			updateDeveloperAppUrl, _ := url.Parse(apiclient.BaseURL)
			updateDeveloperAppUrl.Path = path.Join(updateDeveloperAppUrl.Path, apiclient.GetApigeeOrg(), "developers", developerID, "apps", app.Name, "keys", credential.ConsumerKey)

			updateCred := importCredential{}
			updateCred.ConsumerKey = credential.ConsumerKey
			updateCred.ConsumerSecret = credential.ConsumerSecret
			updateCred.APIProducts = products
			updateCred.Scopes = credential.Scopes

			updateCredJSON, err := json.Marshal(updateCred)
			if err != nil {
				clilog.Error.Println(err)
				return
			}

			_, err = apiclient.HttpClient(updateDeveloperAppUrl.String(), string(updateCredJSON))
			if err != nil {
				return
			}
		} else {
			clilog.Warning.Println("NOTE: apiProducts are not associated with the app")
		}
	}
	clilog.Debug.Printf("Completed entity: %s", app.Name)
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}

func getNewDeveloperId(oldDeveloperId string, developerEntities developers.Appdevelopers) (developerEmail string, newDeveloperId string, err error) {
	if oldDeveloperId == "" {
		return "", "", fmt.Errorf("developer id is null")
	}
	for _, developer := range developerEntities.Developer {
		if oldDeveloperId == developer.DeveloperId {
			newDeveloperId, err = developers.GetDeveloperId(developer.EMail)
			return developer.EMail, newDeveloperId, err
		}
	}
	return "", "", fmt.Errorf("developer not imported into Apigee")
}

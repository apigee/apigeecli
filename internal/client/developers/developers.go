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

package developers

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
)

// Appdevelopers holds a single developer
type Appdeveloper struct {
	EMail       string      `json:"email,omitempty"`
	FirstName   string      `json:"firstName,omitempty"`
	LastName    string      `json:"lastName,omitempty"`
	Attributes  []Attribute `json:"attributes,omitempty"`
	Username    string      `json:"userName,omitempty"`
	DeveloperId string      `json:"developerId,omitempty"`
}

// Appdevelopers hold an array of developers
type Appdevelopers struct {
	Developer []Appdeveloper `json:"developer,omitempty"`
}

// Attribute to used to hold custom attributes for entities
type Attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Create
func Create(email string, firstName string, lastName string, username string, attrs map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	developer := []string{}

	developer = append(developer, "\"email\":\""+email+"\"")
	developer = append(developer, "\"firstName\":\""+firstName+"\"")
	developer = append(developer, "\"lastName\":\""+lastName+"\"")
	developer = append(developer, "\"userName\":\""+username+"\"")

	if len(attrs) != 0 {
		attributes := []string{}
		for key, value := range attrs {
			attributes = append(attributes, "{\"name\":\""+key+"\",\"value\":\""+value+"\"}")
		}
		attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
		developer = append(developer, attributesStr)
	}

	payload := "{" + strings.Join(developer, ",") + "}"
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Delete
func Delete(email string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email)) // since developer emails can have +
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get(email string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email)) // since developer emails can have +
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetDeveloperId
func GetDeveloperId(email string) (developerId string, err error) {
	apiclient.SetClientPrintHttpResponse(false)
	defer apiclient.SetClientPrintHttpResponse(apiclient.GetCmdPrintHttpResponseSetting())
	var developerMap map[string]interface{}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email)) // since developer emails can have +
	respBody, err := apiclient.HttpClient(u.String())
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(respBody, &developerMap)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", developerMap["developerId"]), nil
}

// GetApps
func GetApps(name string, expand bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if expand {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", name, "apps")
		q := u.Query()
		q.Set("expand", "true")
		u.RawQuery = q.Encode()
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", name, "apps")
	}
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// List
func List(count int, expand bool, ids string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers")
	q := u.Query()
	if expand {
		q.Set("expand", "true")
	} else {
		q.Set("expand", "false")
	}
	if count != -1 {
		q.Set("count", strconv.Itoa(count))
	}
	if ids != "" {
		q.Set("ids", ids)
	}

	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Export
func Export() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers")

	q := u.Query()
	q.Set("expand", "true")

	u.RawQuery = q.Encode()
	// don't print to sysout
	apiclient.SetClientPrintHttpResponse(false)
	defer apiclient.SetClientPrintHttpResponse(apiclient.GetCmdPrintHttpResponseSetting())
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Import
func Import(conn int, filePath string) error {
	var pwg sync.WaitGroup
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers")

	entities, err := ReadDevelopersFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	numEntities := len(entities.Developer)
	clilog.Debug.Printf("Found %d developers in the file\n", numEntities)
	clilog.Debug.Printf("Create developers with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	// ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Debug.Printf("Creating batch %d of developers\n", (i + 1))
		go batchImport(u.String(), entities.Developer[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Debug.Printf("Creating remaining %d developers\n", remaining)
		go batchImport(u.String(), entities.Developer[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

func createAsyncDeveloper(url string, dev Appdeveloper, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := json.Marshal(dev)
	if err != nil {
		clilog.Error.Println(err)
		return
	}
	_, err = apiclient.HttpClient(url, string(out))
	if err != nil {
		clilog.Error.Println(err)
		return
	}

	clilog.Debug.Printf("Completed entity: %s", dev.EMail)
}

// batch creates a batch of developers to create
func batchImport(url string, entities []Appdeveloper, pwg *sync.WaitGroup) {
	defer pwg.Done()
	// batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncDeveloper(url, entity, &bwg)
	}
	bwg.Wait()
}

// ReadDevelopersFile
func ReadDevelopersFile(filePath string) (Appdevelopers, error) {
	devs := Appdevelopers{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return devs, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return devs, err
	}

	err = json.Unmarshal(byteValue, &devs)

	if err != nil {
		return devs, err
	}

	return devs, nil
}

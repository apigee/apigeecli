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
	"errors"
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
	Status      *string     `json:"status,omitempty"`
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

// Update
func Update(email string, firstName string, lastName string, username string, status string, attrs map[string]string) (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	devattrs := []Attribute{}

	devRespBody, err := Get(email)
	if err != nil {
		return nil, err
	}
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	d := Appdeveloper{}
	if err = json.Unmarshal(devRespBody, &d); err != nil {
		return nil, err
	}

	if firstName != "" {
		d.FirstName = firstName
	}

	if lastName != "" {
		d.LastName = lastName
	}

	if username != "" {
		d.Username = username
	}

	if len(attrs) > 0 {
		for k, v := range attrs {
			a := Attribute{}
			a.Name = k
			a.Value = v
			devattrs = append(devattrs, a)
		}
		d.Attributes = devattrs
	}

	if status != "" {
		apiclient.ClientPrintHttpResponse.Set(false)
		err = setDeveloperStatus(email, status)
		apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
		if err != nil {
			return nil, err
		}
	}

	reqBody, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email))
	respBody, err = apiclient.HttpClient(u.String(), string(reqBody), "PUT")
	return respBody, err
}

func setDeveloperStatus(email string, action string) (err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email))
	q := u.Query()
	q.Set("action", action)
	u.RawQuery = q.Encode()

	_, err = apiclient.HttpClient(u.String(), "", "POST", "application/octet-stream")
	return err
}

// GetDeveloperId
func GetDeveloperId(email string) (developerId string, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
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
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Import
func Import(conn int, filePath string) error {
	entities, err := ReadDevelopersFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	numEntities := len(entities.Developer)
	clilog.Debug.Printf("Found %d developers in the file\n", numEntities)
	clilog.Debug.Printf("Create developers with %d connections\n", conn)

	jobChan := make(chan Appdeveloper)
	errChan := make(chan error)

	fanOutWg := sync.WaitGroup{}
	fanInWg := sync.WaitGroup{}

	errs := []string{}
	fanInWg.Add(1)
	go func() {
		defer fanInWg.Done()
		for {
			newErr, ok := <-errChan
			if !ok {
				return
			}
			errs = append(errs, newErr.Error())
		}
	}()

	for i := 0; i < conn; i++ {
		fanOutWg.Add(1)
		go createAsyncDeveloper(&fanOutWg, jobChan, errChan)
	}

	for _, entity := range entities.Developer {
		jobChan <- entity
	}
	close(jobChan)
	fanOutWg.Wait()
	close(errChan)
	fanInWg.Wait()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func createAsyncDeveloper(wg *sync.WaitGroup, jobs <-chan Appdeveloper, errs chan<- error) {
	defer wg.Done()
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers")

	for {
		job, ok := <-jobs
		if !ok {
			return
		}
		dev, err := json.Marshal(job)
		if err != nil {
			errs <- err
			continue
		}
		_, err = apiclient.HttpClient(u.String(), string(dev))
		if err != nil {
			errs <- err
			continue
		}
		clilog.Debug.Printf("Completed entity: %s", job.EMail)
	}
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

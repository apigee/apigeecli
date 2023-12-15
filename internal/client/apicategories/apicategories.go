// Copyright 2023 Google LLC
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

package apicategories

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"

	"internal/apiclient"

	"github.com/thedevsaddam/gojsonq"
	"internal/client/sites"
	"internal/clilog"
)

type listapicategories struct {
	Status        string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	RequestID     string `json:"requestId,omitempty"`
	ErrorCode     string `json:"errorCode,omitempty"`
	Data          []data `json:"data,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

type data struct {
	SiteID     string `json:"siteId,omitempty"`
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
}

// Create
func Create(siteid string, name string) (respBody []byte, err error) {
	apicategories := []string{}
	apicategories = append(apicategories, "\"siteId\":"+"\""+siteid+"\"")
	apicategories = append(apicategories, "\"name\":"+"\""+name+"\"")
	payload := "{" + strings.Join(apicategories, ",") + "}"
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apicategories")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Get
func Get(siteid string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apicategories", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// List
func List(siteid string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apicategories")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(siteid string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apicategories", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Update
func Update(siteid string, name string) (respBody []byte, err error) {
	apicategories := []string{}
	apicategories = append(apicategories, "\"siteId\":"+"\""+siteid+"\"")
	apicategories = append(apicategories, "\"name\":"+"\""+name+"\"")
	payload := "{" + strings.Join(apicategories, ",") + "}"
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apicategories")
	respBody, err = apiclient.HttpClient(u.String(), payload, "PATCH")
	return respBody, err
}


// GetByName
func GetByName(siteid string, name string) (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	listRespBytes, err := List(siteid)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch apidocs: %w", err)
	}
	out := gojsonq.New().JSONString(string(listRespBytes)).From("data").Where("name", "eq", name).First()
	if isNil(out) {
		return nil, fmt.Errorf("unable to find category with name %s", name)
	}
	outBytes, err := json.Marshal(out)
	if err != nil {
		return outBytes, err
	}
	return outBytes, nil
}

// from: https://mangatmodi.medium.com/go-check-nil-interface-the-right-way-d142776edef1
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false

// Export
func Export(folder string) (err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	siteids, err := sites.GetSiteIDs()
	if err != nil {
		return err
	}

	for _, siteid := range siteids {
		listRespBytes, err := List(siteid)
		if err != nil {
			return fmt.Errorf("failed to fetch apicategories: %w", err)
		}

		docFileName := fmt.Sprintf("apicategory_%s.json", siteid)
		if err = apiclient.WriteByteArrayToFile(path.Join(folder, docFileName), false, listRespBytes); err != nil {
			return err
		}
	}
	return nil
}

// Import
func Import(siteid string, apicateogyFile string) (err error) {
	errs := []string{}
	l, err := readAPICategoriesFile(apicateogyFile)
	if err != nil {
		return err
	}
	if len(l.Data) < 1 {
		clilog.Warning.Println("No categories found for the siteid")
		return nil
	}

	for _, category := range l.Data {
		_, err = Create(siteid, category.Name)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func readAPICategoriesFile(fileName string) (l listapicategories, err error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return l, err
	}

	defer jsonFile.Close()

	content, err := io.ReadAll(jsonFile)
	if err != nil {
		return l, err
	}
	err = json.Unmarshal(content, &l)
	if err != nil {
		return l, err
	}
	return l, nil
}

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
	"fmt"
	"net/url"
	"path"
	"reflect"
	"strings"

	"internal/apiclient"

	"github.com/thedevsaddam/gojsonq"
)

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
}

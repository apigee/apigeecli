// Copyright 2025 Google LLC
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

package spaces

import (
	"encoding/json"
	"internal/apiclient"
	"net/url"
	"path"
)

type space struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// Create
func Create(name string, displayName string) (respBody []byte, err error) {
	_space := space{
		Name:        name,
		DisplayName: displayName,
	}
	reqBody, err := json.Marshal(_space)
	if err != nil {
		return nil, err
	}
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "spaces")
	respBody, err = apiclient.HttpClient(u.String(), string(reqBody))
	return respBody, err
}

// Update
func Update(name string, displayName string) (respBody []byte, err error) {
	_space := space{
		DisplayName: displayName,
	}
	reqBody, err := json.Marshal(_space)
	if err != nil {
		return nil, err
	}
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "spaces", name)
	respBody, err = apiclient.HttpClient(u.String(), string(reqBody), "PATCH")
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "spaces", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "spaces", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "spaces")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

//TODO: Import, Export

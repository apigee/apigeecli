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

package res

import (
	"fmt"
	"net/url"
	"path"

	"internal/apiclient"
)

// Create
func Create(name string, resPath string, resourceType string) (respBody []byte, err error) {
	if !validate(resourceType) {
		return respBody, fmt.Errorf("invalid resource type")
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "resourcefiles")
	if resourceType != "" {
		q := u.Query()
		q.Set("type", resourceType)
		q.Set("name", name)
		u.RawQuery = q.Encode()
	}
	formParams := map[string]string{
		"file": resPath,
	}
	respBody, err = apiclient.PostHttpOctet(false, u.String(), formParams)
	return respBody, err
}

// Delete
func Delete(name string, resourceType string) (respBody []byte, err error) {
	if !validate(resourceType) {
		return respBody, fmt.Errorf("invalid resource type")
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "resourcefiles", resourceType, name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get(name string, resourceType string) (err error) {
	if !validate(resourceType) {
		return fmt.Errorf("invalid resource type")
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "resourcefiles", resourceType, name)
	err = apiclient.DownloadResource(u.String(), name, resourceType, true)
	return
}

// Update
func Update(name string, resPath string, resourceType string) (err error) {
	if !validate(resourceType) {
		return fmt.Errorf("invalid resource type")
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "resourcefiles", resourceType, name)
	formParams := map[string]string{
		"file": resPath,
	}
	_, err = apiclient.PostHttpOctet(true, u.String(), formParams)
	return
}

// List
func List(resourceType string) (respBody []byte, err error) {
	if resourceType != "" {
		if !validate(resourceType) {
			return respBody, fmt.Errorf("invalid resource type")
		}
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "resourcefiles")

	if resourceType != "" {
		q := u.Query()
		q.Set("type", resourceType)
		u.RawQuery = q.Encode()
	}
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// validate returns true is the resource type is valid
func validate(resType string) bool {
	// validResourceTypes contains a list of valid resources
	validResourceTypes := [7]string{"js", "jsc", "properties", "java", "wsdl", "xsd", "py"}

	for _, n := range validResourceTypes {
		if n == resType {
			return true
		}
	}
	return false
}

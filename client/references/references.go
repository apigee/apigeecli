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

package references

import (
	"encoding/json"
	"net/url"
	"path"
	"strings"

	"github.com/apigee/apigeecli/apiclient"
)

// Create references
func Create(name string, description string, resourceType string, refers string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	reference := []string{}

	reference = append(reference, "\"name\":\""+name+"\"")

	if description != "" {
		reference = append(reference, "\"description\":\""+description+"\"")
	}

	reference = append(reference, "\"resourceType\":\""+resourceType+"\"")
	reference = append(reference, "\"refers\":\""+refers+"\"")

	payload := "{" + strings.Join(reference, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

// Get a reference
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

// Delete a reference
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

// List references
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

// Update references
func Update(name string, description string, resourceType string, refers string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	reference := []string{}

	reference = append(reference, "\"name\":\""+name+"\"")

	if description != "" {
		reference = append(reference, "\"description\":\""+description+"\"")
	}

	if resourceType != "" {
		reference = append(reference, "\"resourceType\":\""+resourceType+"\"")
	}

	if refers != "" {
		reference = append(reference, "\"refers\":\""+refers+"\"")
	}

	payload := "{" + strings.Join(reference, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload, "PUT")
	return respBody, err
}

//Export
func Export(conn int) (payload [][]byte, err error) {
	//TODO: batch exports
	apiclient.SetPrintOutput(false)
	var respBody []byte

	if respBody, err = List(); err != nil {
		return nil, err
	}

	var referencesList []string
	if err := json.Unmarshal(respBody, &referencesList); err != nil {
		return nil, err
	}

	for _, reference := range referencesList {
		if respBody, err = Get(reference); err != nil {
			return nil, err
		}
		payload = append(payload, respBody)
	}

	return payload, nil
}

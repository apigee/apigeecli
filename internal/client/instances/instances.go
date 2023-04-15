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

package instances

import (
	"encoding/json"
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

// Create
func Create(name string, location string, diskEncryptionKeyName string, ipRange string, consumerAcceptList []string) (respBody []byte, err error) {
	instance := []string{}

	instance = append(instance, "\"name\":\""+name+"\"")
	instance = append(instance, "\"location\":\""+location+"\"")

	if ipRange != "" {
		instance = append(instance, "\"ipRange\":\""+ipRange+"\"")
	}

	if diskEncryptionKeyName != "" {
		instance = append(instance, "\"diskEncryptionKeyName\":\""+diskEncryptionKeyName+"\"")
	}

	if len(consumerAcceptList) > 0 {
		builder := new(strings.Builder)
		json.NewEncoder(builder).Encode(consumerAcceptList)
		consumerAcceptListJson := "\"consumerAcceptList\":" + builder.String()
		instance = append(instance, consumerAcceptListJson)
	}

	payload := "{" + strings.Join(instance, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Update
func Update(name string, consumerAcceptList []string) (respBody []byte, err error) {
	instance := []string{}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name)

	if len(consumerAcceptList) > 0 {
		builder := new(strings.Builder)
		json.NewEncoder(builder).Encode(consumerAcceptList)
		consumerAcceptListJson := "\"consumerAcceptList\":" + builder.String()
		instance = append(instance, consumerAcceptListJson)

		q := u.Query()
		q.Set("updateMask", "consumerAcceptList")

		payload := "{" + strings.Join(instance, ",") + "}"
		respBody, err = apiclient.HttpClient(u.String(), payload, "PATCH")
	}
	return respBody, err
}

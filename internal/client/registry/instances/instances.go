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

package instances

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

// Create
func Create(cmekName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "instances")
	q := u.Query()
	q.Set("instanceId", "default")
	u.RawQuery = q.Encode()

	instance := []string{}
	config := []string{}

	name := fmt.Sprintf("projects/%s/locations/%s/instance", apiclient.GetProjectID(), apiclient.GetRegistryRegion())
	instance = append(instance, "\"name\":"+"\""+name+"\"")

	config = append(config, "\"cmekKeyName\":"+"\""+cmekName+"\"")
	configJson := "{" + strings.Join(config, ",") + "}"
	instance = append(instance, "\"config\":"+configJson)

	payload := "{" + strings.Join(instance, ",") + "}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Delete
func Delete() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "instances", "default")
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "instances", "default")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

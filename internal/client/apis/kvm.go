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

package apis

import (
	"net/url"
	"path"
	"strconv"

	"internal/apiclient"
)

// CreateProxyKVM
func CreateProxyKVM(proxyName string, name string, encrypted bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps")
	payload := "{\"name\":\"" + name + "\", \"encrypted\": \"" + strconv.FormatBool(encrypted) + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// DeleteProxyKVM
func DeleteProxyKVM(proxyName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", name, "keyvaluemaps", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// ListProxyKVM
func ListProxyKVM(proxyName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

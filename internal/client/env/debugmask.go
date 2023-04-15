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

package env

import (
	"encoding/json"
	"net/url"
	"path"

	"internal/apiclient"
)

// GetDebug
func GetDebug() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "debugmask")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// SetDebug
func SetDebug(maskConfig string) (respBody []byte, err error) {
	// the following steps will validate json
	m := map[string]string{}
	err = json.Unmarshal([]byte(maskConfig), &m)
	if err != nil {
		return respBody, err
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "debugmask")
	respBody, err = apiclient.HttpClient(u.String(), maskConfig)
	return respBody, err
}

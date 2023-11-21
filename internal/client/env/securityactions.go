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

package env

import (
	"net/url"
	"path"
	"strconv"

	"internal/apiclient"
)

// DisableSecurityAction
func DisableSecurityAction(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "securityActions", name+":disable")
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

// EnableSecurityAction
func EnableSecurityAction(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "securityActions", name+":enable")
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

// GetSecurityAction
func GetSecurityAction(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "securityActions", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListSecurityActions
func ListSecurityActions(pageSize int, pageToken string, filter string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "securityActions")
	q := u.Query()
	if pageSize != -1 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	if filter != "" {
		q.Set("filter", filter)
	}

	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

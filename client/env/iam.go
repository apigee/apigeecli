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
	"net/url"
	"path"

	"github.com/srinandan/apigeecli/apiclient"
)

//GetIAM
func GetIAM() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv()+":getIamPolicy")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//SetIAM
func SetIAM(serviceAccountName string, permission string) (err error) {
	return apiclient.SetIAMServiceAccount(serviceAccountName, permission)
}

//RemoveIAM
func RemoveIAM(serviceAccountName string, role string) (err error) {
	return apiclient.RemoveIAMServiceAccount(serviceAccountName, role)
}

//TestIAM
func TestIAM() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv()+":testIamPermissions")
	payload := "{\"permissions\":[\"apigee.environments.get\"]}"
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

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

package apps

import (
	"net/url"
	"path"
	"strings"

	"github.com/srinandan/apigeecli/apiclient"
)

//CreateKey
func CreateKey(developerEmail string, appID string, consumerKey string, consumerSecret string, apiProducts []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	key := []string{}

	key = append(key, "\"apiProducts\":[\""+getArrayStr(apiProducts)+"\"]")
	key = append(key, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")

	if len(attrs) != 0 {
		attributes := []string{}
		for keyattr, value := range attrs {
			attributes = append(attributes, "{\"name\":\""+keyattr+"\",\"value\":\""+value+"\"}")
		}
		attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
		key = append(key, attributesStr)
	}

	key = append(key, "\"consumerKey\":\""+consumerKey+"\"")
	key = append(key, "\"consumerSecret\":\""+consumerSecret+"\"")

	payload := "{" + strings.Join(key, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerEmail, "apps", appID, "keys")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//DeleteKey
func DeleteKey(developerEmail string, appID string, key string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerEmail, "apps", appID, "keys", key)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//GetKey
func GetKey(developerEmail string, appID string, key string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", developerEmail, "apps", appID, "keys", key)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

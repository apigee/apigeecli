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

package appgroups

import (
	"internal/apiclient"
	"net/url"
	"path"
	"strings"
)

// CreateApp
func CreateApp(name string, expires string, callback string, apiProducts []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	app := []string{}

	app = append(app, "\"name\":\""+name+"\"")

	if len(apiProducts) > 0 {
		app = append(app, "\"apiProducts\":[\""+getArrayStr(apiProducts)+"\"]")
	}

	if callback != "" {
		app = append(app, "\"callbackUrl\":\""+callback+"\"")
	}

	if expires != "" {
		app = append(app, "\"keyExpiresIn\":\""+expires+"\"")
	}

	if len(scopes) > 0 {
		app = append(app, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")
	}

	if len(attrs) != 0 {
		attributes := []string{}
		for key, value := range attrs {
			attributes = append(attributes, "{\"name\":\""+key+"\",\"value\":\""+value+"\"}")
		}
		attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
		app = append(app, attributesStr)
	}

	payload := "{" + strings.Join(app, ",") + "}"
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}

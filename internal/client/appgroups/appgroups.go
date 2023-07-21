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
	"strconv"
	"strings"
)

// Create
func Create(name string, channelUri string, channelId string, displayName string, attrs map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	app := []string{}

	app = append(app, "\"name\":\""+name+"\"")
	if channelUri != "" {
		app = append(app, "\"channelUri\":\""+channelUri+"\"")
	}

	if channelId != "" {
		app = append(app, "\"channelId\":\""+channelId+"\"")
	}

	if displayName != "" {
		app = append(app, "\"displayName\":\""+displayName+"\"")
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

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List(pageSize int, pageToken string, filter string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups")
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

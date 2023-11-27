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

package securityprofiles

import (
	"encoding/json"
	"net/url"
	"path"
	"strconv"
	"strings"

	"internal/apiclient"
)

type secprofile struct {
	Name          string        `json:"name"`
	DisplayName   string        `json:"displayName"`
	Description   string        `json:"description,omitempty"`
	ProfileConfig profileConfig `json:"profileConfig"`
	ScoreConfigs  []scoreConfig `json:"scoreConfigs,omitempty"`
}

type profileConfig struct {
	Categories []category `json:"categories"`
}

type scoreConfig struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ScorePath   string `json:"scorePath,omitempty"`
}

type category struct {
	Abuse         interface{} `json:"abuse,omitempty"`
	Mediation     interface{} `json:"mediation,omitempty"`
	Authorization interface{} `json:"authorization,omitempty"`
	Threat        interface{} `json:"threat,omitempty"`
	Mtls          interface{} `json:"mtls,omitempty"`
	Cors          interface{} `json:"cors,omitempty"`
}

// Create
func Create(name string, content []byte) (respBody []byte, err error) {
	sc := secprofile{}
	if err = json.Unmarshal(content, &sc); err != nil {
		return nil, err
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles")
	q := u.Query()
	q.Set("securityProfileId", name)
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String(), string(content))
	return respBody, err
}

// Attach
func Attach(name string, revision string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name, "environments")
	attachProfile := []string{}
	attachProfile = append(attachProfile, "\"name\":"+"\""+apiclient.GetApigeeEnv()+"\"")
	attachProfile = append(attachProfile, "\"securityProfileRevisionId\":"+"\""+revision+"\"")
	payload := "{" + strings.Join(attachProfile, ",") + "}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Detach
func Detach(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles",
		name, "environments", apiclient.GetApigeeEnv())
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListVersions
func ListVersions(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name+":listRevisions")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// List
func List(pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles")
	q := u.Query()
	if pageSize != -1 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}

	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

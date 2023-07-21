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
	"encoding/json"
	"fmt"
	"internal/apiclient"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type appgroupapps struct {
	AppGroupApps  []appgroupapp `json:"appGroupApps,omitempty"`
	NextPageToken string        `json:"nextPageToken,omitempty"`
}

type appgroupapp struct {
	AppId          string       `json:"appId,omitempty"`
	Name           string       `json:"name,omitempty"`
	Status         string       `json:"status,omitempty"`
	AppGroup       string       `json:"appGroup,omitempty"`
	CreatedAt      string       `json:"createdAt,omitempty"`
	LastModifiedAt string       `json:"lastModifiedAt,omitempty"`
	Credentials    []credential `json:"credentials,omitempty"`
}

type credential struct {
	APIProducts    []apiProduct `json:"apiProducts,omitempty"`
	ConsumerKey    string       `json:"consumerKey,omitempty"`
	ConsumerSecret string       `json:"consumerSecret,omitempty"`
	ExpiresAt      string       `json:"expiresAt,omitempty"`
	Status         string       `json:"status,omitempty"`
	Scopes         []string     `json:"scopes,omitempty"`
}

type apiProduct struct {
	ApiProduct string `json:"apiproduct,omitempty"`
	Status     string `json:"status,omitempty"`
}

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

// DeleteApp
func DeleteApp(name string, appName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetApp
func GetApp(name string, appName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// ListApps
func ListApps(name string, pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps")
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

// ManageApp
func ManageApp(name string, appName string, action string) (respBody []byte, err error) {
	if action != "revoke" && action != "approve" {
		return nil, fmt.Errorf("invalid action. action must be revoke or approve")
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName)
	q := u.Query()
	q.Set("action", action)
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), "", "POST", "application/octet-stream")
	return respBody, err
}

// ExportApps
func ExportApps(name string) (respBody []byte, err error) {
	// don't print to sysout
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	pageToken := ""
	applist := appgroupapps{}

	for {
		a := appgroupapps{}
		listRespBytes, err := ListApps(name, maxPageSize, pageToken)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(listRespBytes, &a)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall: %w", err)
		}
		applist.AppGroupApps = append(applist.AppGroupApps, a.AppGroupApps...)
		pageToken = a.NextPageToken
		if a.NextPageToken == "" {
			break
		}
	}

	respBody, err = json.Marshal(applist.AppGroupApps)

	return respBody, err
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}

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

package apis

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"internal/apiclient"
)

// Create
func Create(apiID string, name string, displayName string, description string,
	availability string, recommendedVersion string, recommendedDeployment string,
	labels map[string]string, annotations map[string]string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis")
	q := u.Query()
	q.Set("apiId", apiID)
	u.RawQuery = q.Encode()

	apiContent := []string{}
	apiContent = append(apiContent, "\"name\":"+"\""+name+"\"")
	if displayName != "" {
		apiContent = append(apiContent, "\"displayName\":"+"\""+displayName+"\"")
	}
	if description != "" {
		apiContent = append(apiContent, "\"description\":"+"\""+description+"\"")
	}
	if availability != "" {
		apiContent = append(apiContent, "\"availability\":"+"\""+availability+"\"")
	}
	if recommendedVersion != "" {
		apiContent = append(apiContent, "\"recommendedVersion\":"+"\""+recommendedVersion+"\"")
	}
	if recommendedDeployment != "" {
		apiContent = append(apiContent, "\"recommendedDeployment\":"+"\""+recommendedDeployment+"\"")
	}
	if len(labels) > 0 {
		l := []string{}
		for key, value := range labels {
			l = append(l, "\""+key+"\":\""+value+"\"")
		}
		labelStr := "\"labels\":{" + strings.Join(l, ",") + "}"
		apiContent = append(apiContent, labelStr)
	}

	if len(annotations) > 0 {
		a := []string{}
		for key, value := range annotations {
			a = append(a, "\""+key+"\":\""+value+"\"")
		}
		annotationStr := "\"annotations\":{" + strings.Join(a, ",") + "}"
		apiContent = append(apiContent, annotationStr)
	}
	payload := "{" + strings.Join(apiContent, ",") + "}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List(pageSize int, pageToken string, filter string, orderBy string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis")
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
	if orderBy != "" {
		q.Set("orderBy", orderBy)
	}

	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

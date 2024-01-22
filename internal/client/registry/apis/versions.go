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

// CreateVersion
func CreateVersion(apiVersionId string, name string, displayName string,
	description string, state string, labels map[string]string,
	annotations map[string]string, primarySpec string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "versions")
	apiVersionContent := []string{}
	apiVersionContent = append(apiVersionContent, "\"name\":"+"\""+name+"\"")
	if displayName != "" {
		apiVersionContent = append(apiVersionContent, "\"displayName\":"+"\""+displayName+"\"")
	}
	if description != "" {
		apiVersionContent = append(apiVersionContent, "\"description\":"+"\""+description+"\"")
	}
	if state != "" {
		apiVersionContent = append(apiVersionContent, "\"state\":"+"\""+state+"\"")
	}
	if len(labels) > 0 {
		l := []string{}
		for key, value := range labels {
			l = append(l, "\""+key+"\":\""+value+"\"")
		}
		labelStr := "\"labels\":{" + strings.Join(l, ",") + "}"
		apiVersionContent = append(apiVersionContent, labelStr)
	}

	if len(annotations) > 0 {
		a := []string{}
		for key, value := range annotations {
			a = append(a, "\""+key+"\":\""+value+"\"")
		}
		annotationStr := "\"annotations\":{" + strings.Join(a, ",") + "}"
		apiVersionContent = append(apiVersionContent, annotationStr)
	}
	payload := "{" + strings.Join(apiVersionContent, ",") + "}"

	q := u.Query()
	q.Set("apiVersionId", apiVersionId)
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// GetVersion
func GetVersion(name string, version string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "versions", version)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// DeleteVersion
func DeleteVersion(name string, version string, force bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "versions", version)
	if force {
		q := u.Query()
		q.Set("force", strconv.FormatBool(force))
		u.RawQuery = q.Encode()
	}
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// ListVersions
func ListVersions(name string, pageSize int, pageToken string,
	filter string, orderBy string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "versions")
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

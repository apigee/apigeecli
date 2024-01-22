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

package common

import (
	"net/url"
	"strconv"
	"strings"

	"internal/apiclient"
)

// GetArtifactPayload
func GetArtifactPayload(name string, contents string, labels map[string]string, annotations map[string]string) string {
	artifact := []string{}
	artifact = append(artifact, "\"name\":"+"\""+name+"\"")
	artifact = append(artifact, "\"contents\":"+"\""+contents+"\"")

	if len(labels) > 0 {
		l := []string{}
		for key, value := range labels {
			l = append(l, "\""+key+"\":\""+value+"\"")
		}
		labelStr := "\"labels\":{" + strings.Join(l, ",") + "}"
		artifact = append(artifact, labelStr)
	}

	if len(annotations) > 0 {
		a := []string{}
		for key, value := range annotations {
			a = append(a, "\""+key+"\":\""+value+"\"")
		}
		annotationStr := "\"annotations\":{" + strings.Join(a, ",") + "}"
		artifact = append(artifact, annotationStr)
	}

	payload := "{" + strings.Join(artifact, ",") + "}"
	return payload
}

// ListArtifacts
func ListArtifacts(u *url.URL, pageSize int,
	pageToken string, filter string, orderBy string,
) (respBody []byte, err error) {
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

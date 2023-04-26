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

package apis

import (
	"bytes"
	"fmt"
	"net/url"
	"path"
	"strconv"

	"internal/apiclient"
)

func getFilterStr(filter map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range filter {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

// CreateTraceSession
func CreateTraceSession(name string, revision int, filter map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"apis", name, "revisions", strconv.Itoa(revision), "debugsessions")
	q := u.Query()
	q.Set("timeout", "567")
	u.RawQuery = q.Encode()

	var payload string
	if len(filter) != 0 {
		payload = "{\"filter\":" + getFilterStr(filter) + "}"
	} else {
		payload = "{}"
	}
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// GetTraceSession
func GetTraceSession(name string, revision int, sessionID string, messageID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if messageID == "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
			"apis", name, "revisions", strconv.Itoa(revision), "debugsessions", sessionID, "data")
		q := u.Query()
		q.Set("limit", "20")
		u.RawQuery = q.Encode()
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
			"apis", name, "revisions", strconv.Itoa(revision), "debugsessions", sessionID, "data", messageID)
	}

	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListTracceSession
func ListTracceSession(name string, revision int) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"apis", name, "revisions", strconv.Itoa(revision), "debugsessions")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

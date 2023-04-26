// Copyright 2021 Google LLC
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

package eptattachment

import (
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

// Create
func Create(name string, serviceAttachment string, location string) (respBody []byte, err error) {
	endpointAttachment := []string{}

	endpointAttachment = append(endpointAttachment, "\"serviceAttachment\":\""+serviceAttachment+"\"")
	endpointAttachment = append(endpointAttachment, "\"location\":\""+location+"\"")

	payload := "{" + strings.Join(endpointAttachment, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "endpointAttachments")

	q := u.Query()
	q.Set("endpointAttachmentId", name)
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "endpointAttachments", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "endpointAttachments", name)
	respBody, err = apiclient.HttpClient(u.String(), "DELETE")
	return respBody, err
}

// List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "endpointAttachments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

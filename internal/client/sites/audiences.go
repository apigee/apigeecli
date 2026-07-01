// Copyright 2024 Google LLC
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

package sites

import (
	"internal/apiclient"
	"net/url"
	"path"
)

// ListAudiences returns all audiences for the given portal site.
func ListAudiences(siteID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteID, "audiences")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetAudience returns a single audience by ID.
func GetAudience(siteID string, audienceID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteID, "audiences", audienceID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// CreateAudience creates a new portal audience from a JSON payload.
func CreateAudience(siteID string, audienceID string, payload []byte) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteID, "audiences")
	if audienceID != "" {
		q := u.Query()
		q.Set("audienceId", audienceID)
		u.RawQuery = q.Encode()
	}
	respBody, err = apiclient.HttpClient(u.String(), string(payload))
	return respBody, err
}

// UpdateAudience updates an existing portal audience (full replace).
func UpdateAudience(siteID string, audienceID string, payload []byte) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteID, "audiences", audienceID)
	respBody, err = apiclient.HttpClient(u.String(), string(payload), "PUT")
	return respBody, err
}

// DeleteAudience deletes a portal audience.
func DeleteAudience(siteID string, audienceID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteID, "audiences", audienceID)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// ManageDeveloperAudience adds or removes a developer from a portal audience.
// action must be "add" or "remove".
func ManageDeveloperAudience(siteID string, audienceID string, developerEmail string, action string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteID,
		"audiences", audienceID, "developers:"+action)
	payload := `{"developers":["` + developerEmail + `"]}`
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

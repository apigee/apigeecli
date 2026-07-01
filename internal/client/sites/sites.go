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

package sites

import (
	"internal/apiclient"
	"net/url"
	"path"

	"github.com/thedevsaddam/gojsonq"
)

// List returns all integrated developer portals for the org.
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Get returns a single integrated developer portal by site ID.
func Get(siteID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Create creates a new integrated developer portal from a JSON payload.
// siteID is the desired portal ID; it may be passed as a query param.
func Create(siteID string, payload []byte) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites")
	if siteID != "" {
		q := u.Query()
		q.Set("siteId", siteID)
		u.RawQuery = q.Encode()
	}
	respBody, err = apiclient.HttpClient(u.String(), string(payload))
	return respBody, err
}

// Update replaces an integrated developer portal configuration.
func Update(siteID string, payload []byte) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteID)
	respBody, err = apiclient.HttpClient(u.String(), string(payload), "PUT")
	return respBody, err
}

// Delete removes an integrated developer portal.
func Delete(siteID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteID)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// GetSiteIDs
func GetSiteIDs() (siteIDs []string, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	respBody, err := List()
	if err != nil {
		return nil, err
	}
	jq := gojsonq.New().JSONString(string(respBody))
	ids := jq.From("data").Pluck("id").([]interface{})
	siteIDs = make([]string, len(ids))
	for k, v := range ids {
		siteIDs[k] = v.(string)
	}
	return siteIDs, nil
}

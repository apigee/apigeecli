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
	"fmt"
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

// CreateKey
func CreateKey(name string, appName string, consumerKey string, consumerSecret string, expiresInSeconds string, apiProducts []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	key := []string{}

	key = append(key, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")

	if len(attrs) != 0 {
		attributes := []string{}
		for keyattr, value := range attrs {
			attributes = append(attributes, "{\"name\":\""+keyattr+"\",\"value\":\""+value+"\"}")
		}
		attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
		key = append(key, attributesStr)
	}

	if consumerKey != "" {
		key = append(key, "\"consumerKey\":\""+consumerKey+"\"")
		key = append(key, "\"consumerSecret\":\""+consumerSecret+"\"")
	}

	payload := "{" + strings.Join(key, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName, "keys")

	if len(apiProducts) > 0 {
		apiclient.ClientPrintHttpResponse.Set(false)
	}
	respBody, err = apiclient.HttpClient(u.String(), payload)

	if err != nil {
		return respBody, err
	}

	// since the API does not support adding products when creating a key, use a second API call to add products
	if len(apiProducts) > 0 {
		apiclient.ClientPrintHttpResponse.Set(false)
		respBody, err = UpdateKeyProducts(name, appName, consumerKey, apiProducts)
		apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	}

	return respBody, err
}

// GetKey
func GetKey(name string, appName string, key string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName, "keys", key)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// DeleteKey
func DeleteKey(name string, appName string, key string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName, "keys", key)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// ManageKey
func ManageKey(name string, appName string, consumerKey string, action string) (respBody []byte, err error) {
	if action != "revoke" && action != "approve" {
		return nil, fmt.Errorf("invalid action. action must be revoke or approve")
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName, "keys", consumerKey)
	q := u.Query()
	q.Set("action", action)
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String(), "", "POST", "application/octet-stream")

	return respBody, err
}

// UpdateKeyProducts
func UpdateKeyProducts(name string, appName string, consumerKey string, apiProducts []string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	key := []string{}
	key = append(key, "\"apiProducts\":[\""+getArrayStr(apiProducts)+"\"]")

	payload := "{" + strings.Join(key, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName, "keys", consumerKey)
	respBody, err = apiclient.HttpClient(u.String(), payload)

	return respBody, err
}

// DeleteKeyProduct
func DeleteKeyProduct(name string, appName string, consumerKey string, productName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName, "keys", consumerKey, "apiproducts", productName)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

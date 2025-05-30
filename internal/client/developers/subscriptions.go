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

package developers

import (
	"internal/apiclient"
	"net/url"
	"path"
	"strings"
)

func CreateSubscription(email string, apiproduct string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())

	subscription := []string{}
	subscription = append(subscription, "\"apiproduct\":\""+apiproduct+"\"")

	payload := "{" + strings.Join(subscription, ",") + "}"
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email), "subscriptions")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// ExpireSubscriptions
func ExpireSubscriptions(email string, subscription string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email), "subscriptions", subscription, ":expire") // since developer emails can have +
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

// GetSubscriptions
func GetSubscriptions(email string, subscription string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email), "subscriptions", subscription) // since developer emails can have +
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListSubscriptions
func ListSubscriptions(email string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email), "subscriptions") // since developer emails can have +
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ExportSubscriptions
func ExportSubscriptions(email string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email), "subscriptions")

	// don't print to sysout
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

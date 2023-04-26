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

package instances

import (
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

// ReserveNatIP
func ReserveNatIP(name string, natid string) (respBody []byte, err error) {
	reserveNat := []string{}
	reserveNat = append(reserveNat, "\"name\":\""+natid+"\"")
	payload := "{" + strings.Join(reserveNat, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name, "natAddresses")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// ActivateNatIP
func ActivateNatIP(name string, natid string) (respBody []byte, err error) {
	payload := "{}"
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name, "natAddresses", natid, ":activate")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// DeleteNatIP
func DeleteNatIP(name string, natid string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name, "natAddresses", natid)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// ListNatIPs
func ListNatIPs(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name, "natAddresses")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

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

package sync

import (
	"encoding/json"
	"net/url"
	"path"
	"strings"

	"github.com/srinandan/apigeecli/apiclient"
)

type iAMIdentities struct {
	Identities []string `json:"identities,omitempty"`
}

func validate(i string) string {
	if strings.Contains(i, "serviceAccount:") {
		return i
	}
	return "serviceAccount:" + i
}

func Get() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":getSyncAuthorization")
	respBody, err = apiclient.HttpClient(true, u.String(), "")
	return respBody, err
}

func Reset() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":setSyncAuthorization")
	payload := "{\"identities\":[]}"
	respBody, err = apiclient.HttpClient(true, u.String(), payload)
	return respBody, err
}

func Set(identity string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":setSyncAuthorization")
	identity = validate(identity)
	identities := iAMIdentities{}
	identities.Identities = append(identities.Identities, identity)
	payload, _ := json.Marshal(&identities)
	respBody, err = apiclient.HttpClient(true, u.String(), string(payload))
	return respBody, err
}

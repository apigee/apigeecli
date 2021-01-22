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
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/srinandan/apigeecli/apiclient"
)

type iAMIdentities struct {
	Identities []string `json:"identities,omitempty"`
}

type syncResponse struct {
	Identities []string `json:"identities,omitempty"`
	Etag       string   `json:"etag,omitempty"`
}

func validate(i string) string {
	if strings.Contains(i, "serviceAccount:") {
		return i
	}
	return "serviceAccount:" + i
}

//Get
func Get() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":getSyncAuthorization")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "")
	return respBody, err
}

//Reset
func Reset() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":setSyncAuthorization")
	payload := "{\"identities\":[]}"
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//Set
func Set(identity string) (respBody []byte, err error) {

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":getSyncAuthorization")
	apiclient.SetPrintOutput(false)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "")
	if err != nil {
		return respBody, err
	}

	response := syncResponse{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return respBody, err
	}

	identity = validate(identity)

	for _, setIdentity := range response.Identities {
		if identity == setIdentity {
			return respBody, fmt.Errorf("identity %s already set", identity)
		}
	}

	response.Identities = append(response.Identities, identity)

	identities := iAMIdentities{}
	identities.Identities = response.Identities
	payload, err := json.Marshal(&identities)
	if err != nil {
		return respBody, err
	}

	apiclient.SetPrintOutput(true)
	u, _ = url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":setSyncAuthorization")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(payload))

	return respBody, err
}

//SetList
func SetList(identities []string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":getSyncAuthorization")
	apiclient.SetPrintOutput(false)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "")
	if err != nil {
		return respBody, err
	}

	response := syncResponse{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return respBody, err
	}

	for count := 0; count < len(identities); count++ {
		identities[count] = validate(identities[count])
	}

	response.Identities = append(response.Identities, identities...)

	iamidentities := iAMIdentities{}
	iamidentities.Identities = response.Identities
	payload, err := json.Marshal(&iamidentities)
	if err != nil {
		return respBody, err
	}

	apiclient.SetPrintOutput(true)
	u, _ = url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":setSyncAuthorization")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(payload))

	return respBody, err
}

//Remove
func Remove(identity string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":getSyncAuthorization")
	apiclient.SetPrintOutput(false)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "")
	if err != nil {
		return respBody, err
	}

	response := syncResponse{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return respBody, err
	}

	identity = validate(identity)
	found := false

	numIdentities := len(response.Identities)
	if numIdentities < 1 {
		return respBody, fmt.Errorf("identity %s not found", identity)
	}

	for i, setIdentity := range response.Identities {
		if identity == setIdentity {
			response.Identities[i] = response.Identities[numIdentities-1]
			response.Identities[numIdentities-1] = ""
			response.Identities = response.Identities[:numIdentities-1]
			found = true
		}
	}

	if !found {
		return respBody, fmt.Errorf("identity %s not found", identity)
	}

	identities := iAMIdentities{}
	identities.Identities = response.Identities
	payload, err := json.Marshal(&identities)
	if err != nil {
		return respBody, err
	}

	apiclient.SetPrintOutput(true)
	u, _ = url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":setSyncAuthorization")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(payload))

	return respBody, err
}

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

package kvm

import (
	"encoding/json"
	"net/url"
	"path"
	"strconv"

	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/clilog"
)

type keyvalueentry struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type keyvalueentries struct {
	KeyValueEntries []keyvalueentry `json:"keyValueEntries,omitempty"`
	NextPageToken   string          `json:"nextPageToken,omitempty"`
}

//CreateEntry
func CreateEntry(proxyName string, mapName string, keyName string, value string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", mapName, "entries")
	} else if proxyName != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps", mapName, "entries")
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "keyvaluemaps", mapName, "entries")
	}
	payload := "{\"name\":\"" + keyName + "\",\"value\":\"" + value + "\"}"
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//DeleteEntry
func DeleteEntry(proxyName string, mapName string, keyName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", mapName, "entries", keyName)
	} else if proxyName != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps", mapName, "entries", keyName)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "keyvaluemaps", mapName, "entries", keyName)
	}
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//GetEntry
func GetEntry(proxyName string, mapName string, keyName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if apiclient.GetApigeeEnv() != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", mapName, "entries", keyName)
	} else if proxyName != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps", mapName, "entries", keyName)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "keyvaluemaps", mapName, "entries", keyName)
	}
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ListEntries
func ListEntries(proxyName string, mapName string, pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	if pageToken != "" || pageSize != -1 {
		q := u.Query()
		if pageToken != "" {
			q.Set("page_token", pageToken)
		}
		if pageSize != -1 {
			q.Set("page_size", strconv.Itoa(pageSize))
		}
		u.RawQuery = q.Encode()
	}

	if apiclient.GetApigeeEnv() != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", mapName, "entries")
	} else if proxyName != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps", mapName, "entries")
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "keyvaluemaps", mapName, "entries")
	}
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ExportEntries
func ExportEntries(proxyName string, mapName string) (payload [][]byte, err error) {
	var respBody []byte
	count := 1

	apiclient.SetPrintOutput(false)
	if respBody, err = ListEntries(proxyName, mapName, -1, ""); err != nil {
		return nil, err
	}

	clilog.Info.Printf("Exporting batch 1 of KVM entries for map %s\n", mapName)
	payload = append(payload, respBody)

	var keyValueEntries = keyvalueentries{}
	if err = json.Unmarshal(respBody, &keyValueEntries); err != nil {
		return nil, err
	}

	for keyValueEntries.NextPageToken != "" {
		if respBody, err = ListEntries(proxyName, mapName, -1, keyValueEntries.NextPageToken); err != nil {
			return nil, err
		}

		if err = json.Unmarshal(respBody, &keyValueEntries); err != nil {
			return nil, err
		}
		count++
		clilog.Info.Printf("Exporting batch %d of KVM entries for map %s\n", count, mapName)
		payload = append(payload, respBody)
	}

	apiclient.SetPrintOutput(true)
	return payload, nil
}

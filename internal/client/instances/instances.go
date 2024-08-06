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
	"encoding/json"
	"net/url"
	"path"
	"strings"
	"time"

	"internal/apiclient"
	"internal/clilog"
)

// Create
func Create(name string, location string, diskEncryptionKeyName string, ipRange string, consumerAcceptList []string) (respBody []byte, err error) {
	instance := []string{}

	instance = append(instance, "\"name\":\""+name+"\"")
	instance = append(instance, "\"location\":\""+location+"\"")

	if ipRange != "" {
		instance = append(instance, "\"ipRange\":\""+ipRange+"\"")
	}

	if diskEncryptionKeyName != "" {
		instance = append(instance, "\"diskEncryptionKeyName\":\""+diskEncryptionKeyName+"\"")
	}

	if len(consumerAcceptList) > 0 {
		builder := new(strings.Builder)
		json.NewEncoder(builder).Encode(consumerAcceptList)
		consumerAcceptListJson := "\"consumerAcceptList\":" + builder.String()
		instance = append(instance, consumerAcceptListJson)
	}

	payload := "{" + strings.Join(instance, ",") + "}"

	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Update
func Update(name string, consumerAcceptList []string) (respBody []byte, err error) {
	instance := []string{}

	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name)

	if len(consumerAcceptList) > 0 {
		builder := new(strings.Builder)
		json.NewEncoder(builder).Encode(consumerAcceptList)
		consumerAcceptListJson := "\"consumerAcceptList\":" + builder.String()
		instance = append(instance, consumerAcceptListJson)

		q := u.Query()
		q.Set("updateMask", "consumerAcceptList")

		payload := "{" + strings.Join(instance, ",") + "}"
		respBody, err = apiclient.HttpClient(u.String(), payload, "PATCH")
	}
	return respBody, err
}

// Wait
func Wait(name string) error {
	var err error

	clilog.Info.Printf("Checking creation status in %d seconds\n", interval)

	apiclient.DisableCmdPrintHttpResponse()

	stop := apiclient.Every(interval*time.Second, func(time.Time) bool {
		var respBody []byte
		respMap := make(map[string]interface{})
		if respBody, err = Get(name); err != nil {
			clilog.Error.Printf("Error fetching env status: %v", err)
			return false
		}

		if err = json.Unmarshal(respBody, &respMap); err != nil {
			return true
		}

		switch respMap["state"] {
		case "PROGRESSING":
			clilog.Info.Printf("Instance creation status is: %s. Waiting %d seconds.\n", respMap["state"], interval)
			return true
		case "ACTIVE":
			clilog.Info.Println("Instance creation completed with status: ", respMap["state"])
		default:
			clilog.Info.Println("Instance creation failed with status: ", respMap["state"])
		}

		return false
	})

	<-stop

	return err
}

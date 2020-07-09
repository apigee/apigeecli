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

package envgroups

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/srinandan/apigeecli/apiclient"
)

//Create
func Create(name string, hostnames []string) (respBody []byte, err error) {

	envgroup := []string{}

	envgroup = append(envgroup, "\"name\":\""+name+"\"")
	envgroup = append(envgroup, "\"hostnames\":[\""+getArrayStr(hostnames)+"\"]")

	payload := "{" + strings.Join(envgroup, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "envgroups")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "envgroups", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "envgroups", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "DELETE")
	return respBody, err
}

//List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "envgroups")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//PatchHosts
func PatchHosts(name string, hostnames []string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "envgroups", name)
	q := u.Query()
	q.Set("updateMask", "hostnames")
	u.RawQuery = q.Encode()

	envgroup := []string{}
	envgroup = append(envgroup, "\"hostnames\":[\""+getArrayStr(hostnames)+"\"]")
	payload := "{" + strings.Join(envgroup, ",") + "}"

	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload, "PATCH", "application/merge-patch+json")
	return respBody, err
}

//Attach
func Attach(name string, environment string) (respBody []byte, err error) {

	envgroup := []string{}
	envgroup = append(envgroup, "\"environment\":\""+environment+"\"")
	payload := "{" + strings.Join(envgroup, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "envgroups", name, "attachments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//DetachEnvironment
func DetachEnvironment(name string, environment string) (respBody []byte, err error) {

	type attachment struct {
		Name        string `json:"name,omitempty"`
		Environment string `json:"environment,omitempty"`
		CreatedAt   string `json:"createdAt,omitempty"`
	}

	type attachments struct {
		Attachment []attachment `json:"environmentGroupAttachments,omitempty"`
	}

	envGroupAttachments := attachments{}

	apiclient.SetPrintOutput(false)
	if respBody, err = ListAttach(name); err != nil {
		return nil, err
	}
	apiclient.SetPrintOutput(true)

	if err := json.Unmarshal(respBody, &envGroupAttachments); err != nil {
		return nil, err
	}

	for _, envGroupAttachment := range envGroupAttachments.Attachment {
		if envGroupAttachment.Environment == environment {
			respBody, err = Detach(name, envGroupAttachment.Name)
			return respBody, err
		}
	}

	return nil, fmt.Errorf("did not find environment %s in envgroup %s", environment, name)
}

//Detach
func Detach(name string, attachment string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "envgroups", name, "attachments", attachment)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//ListAttach
func ListAttach(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "envgroups", name, "attachments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}

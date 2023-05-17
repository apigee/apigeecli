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
	"fmt"
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

// Attach
func Attach(name string, environment string) (respBody []byte, err error) {
	envgroup := []string{}
	envgroup = append(envgroup, "\"environment\":\""+environment+"\"")
	payload := "{" + strings.Join(envgroup, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name, "attachments")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// DetachEnv
func DetachEnv(instance string) (respBody []byte, err error) {
	var attachmentName string
	u, _ := url.Parse(apiclient.BaseURL)

	if attachmentName, err = getAttachmentName(instance); err != nil {
		return nil, err
	}

	if attachmentName == "" {
		return nil, fmt.Errorf("The environment %s, does not appear to be attached to the instance %s", apiclient.GetApigeeEnv(), instance)
	}

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", instance, "attachments", attachmentName)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// func GetEnv
func GetEnv(instance string) (respBody []byte, err error) {
	var attachmentName string
	u, _ := url.Parse(apiclient.BaseURL)

	if attachmentName, err = getAttachmentName(instance); err != nil {
		return nil, err
	}

	if attachmentName == "" {
		return nil, fmt.Errorf("The environment %s, does not appear to be attached to the instance %s", apiclient.GetApigeeEnv(), instance)
	}

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", instance, "attachments", attachmentName)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Detach
func Detach(name string, instanceName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", instanceName, "attachments", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// GetAttach
func GetAttach(name string, instanceName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", instanceName, "attachments", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListAttach
func ListAttach(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "instances", name, "attachments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// getAttachmentName
func getAttachmentName(instance string) (attachmentName string, err error) {
	type instanceAttachment struct {
		Name        string `json:"name,omitempty"`
		Environment string `json:"environment,omitempty"`
		CreatedAt   string `json:"createdAt,omitempty"`
	}

	type instanceAttachments struct {
		Attachments []instanceAttachment `json:"attachments,omitempty"`
	}

	instAttach := instanceAttachments{}

	apiclient.ClientPrintHttpResponse.Set(false)
	listAttachments, err := ListAttach(instance)
	if err != nil {
		return "", err
	}
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	err = json.Unmarshal(listAttachments, &instAttach)
	if err != nil {
		return "", err
	}

	if len(instAttach.Attachments) < 1 {
		return "", fmt.Errorf("no environments attached to the instance")
	}

	for _, attachedEnv := range instAttach.Attachments {
		if attachedEnv.Environment == apiclient.GetApigeeEnv() {
			attachmentName = attachedEnv.Name
			break
		}
	}
	return attachmentName, nil
}

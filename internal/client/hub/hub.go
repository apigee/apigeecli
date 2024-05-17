// Copyright 2024 Google LLC
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

package hub

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

func CreateInstance(apiHubInstanceId string, cmekName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apiHubInstances")

	if apiHubInstanceId != "" {
		q := u.Query()
		q.Set("apiHubInstanceId", apiHubInstanceId)
		u.RawQuery = q.Encode()
	}

	instance := []string{}
	config := []string{}

	name := fmt.Sprintf("projects/%s/locations/%s/instance", apiclient.GetProjectID(), apiclient.GetRegistryRegion())
	instance = append(instance, "\"name\":"+"\""+name+"\"")

	config = append(config, "\"cmekKeyName\":"+"\""+cmekName+"\"")
	configJson := "{" + strings.Join(config, ",") + "}"
	instance = append(instance, "\"config\":"+configJson)

	payload := "{" + strings.Join(instance, ",") + "}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func GetInstance(apiHubInstanceId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apiHubInstances", apiHubInstanceId)

	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func RegisterHostProject(registrationId string, gcpProjectId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "hostProjectRegistrations")

	q := u.Query()
	q.Set("hostProjectRegistrationId", registrationId)
	u.RawQuery = q.Encode()

	payload := fmt.Sprintf("{\"gcpProject\":\"%s\"}", gcpProjectId)
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func CreateRuntimeProjectAttachment(runtimeProjectAttachmentId string, runtimeProject string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "runtimeProjectAttachments")

	q := u.Query()
	q.Set("runtimeProjectAttachmentId", runtimeProjectAttachmentId)
	u.RawQuery = q.Encode()

	payload := fmt.Sprintf("{\"runtimeProject\":\"%s\"}", runtimeProject)
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func DeleteRuntimeProjectAttachment(runtimeProjectAttachmentId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "runtimeProjectAttachments", runtimeProjectAttachmentId)

	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func GetRuntimeProjectAttachment(runtimeProjectAttachmentId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "runtimeProjectAttachments", runtimeProjectAttachmentId)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func ListRuntimeProjectAttachment(runtimeProjectAttachmentId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "runtimeProjectAttachments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

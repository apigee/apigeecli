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
	"encoding/base64"
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

func ListHostProjects() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "hostProjectRegistrations")
	respBody, err = apiclient.HttpClient(u.String())
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

func CreateApi(apiID string, contents byte[]) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis")

	q := u.Query()
	q.Set("apiId", apiID)
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), contents)
	return respBody, err
}

func DeleteApi(apiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID)

	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func GetApi(apiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID)

	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func ListApi(filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID)
	q := u.Query()
	if filter != "" {
		q.Set("filter", filter)
	}
	if pageSize != -1 {
		q.Set("pageSize", pageSize)
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func CreateApiVersion(versionID string, apiID string, contents byte[]) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID)
	q := u.Query()
	if filter != "" {
		q.Set("versionId", filter)
	}
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String(), contents)
	return respBody, err
}

func GetApiVersion(versionID, apiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteApiVersion(versionID, apiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func ListApiVersions(apiID string, filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions")
	q := u.Query()
	if filter != "" {
		q.Set("filter", filter)
	}
	if pageSize != -1 {
		q.Set("pageSize", pageSize)
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func GetApiVersionsDefinitions(apiID string, versionID string, definition string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "definitions", definition)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

//TODO: specType
func CreateApiVersionsSpec(apiID string, versionID string, specID string, displayName string,
	contents []byte, mimeType string, attributes map[string]string, sourceURI string, documentation string) (respBody []byte, err error) {

	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs")
	q := u.Query()
	if filter != "" {
		q.Set("specId", specID)
	}
	u.RawQuery = q.Encode()

	spec := []string{}
	contentPayload := ""
	spec = append(spec, "\"displayName\":"+"\""+displayName+"\"")
	if sourceURI != "" {
		specContent = append(specContent, "\"sourceUri\":"+"\""+sourceURI+"\"")
	}
	if documentation != "" {
		specContent = append(specContent, "\"documentation\":"+"\""+documentation+"\"")
	}
	if contents != nil {
		contentStr := []string{}
		mime := ""
		if mimeType != "" {
			if strings.Contains(mimeType, "yaml") || strings.Contains(mimeType, "yml") {
				mime = "application/yaml"
			} else if strings.Contains(mimeType, "json") {
				mime = "application/json"
			} else if strings.Contains(mimeType, "wsdl") {
				mime = "application/wsdl"
			} else {
				mime = "application/text"
			}
		} else {
			mime = "application/text"
		}
		contentStr = append(contentStr, "\"mimeType\":"+"\""+mime+"\"")
		encContent := base64.StdEncoding.EncodeToString(contents)
		contentStr = append(contentStr, "\"contents\":"+"\""+encContent+"\"")
		contentPayload := "{" + strings.Join(contentStr, ",") + "}"
		specContent = append(specContent, "\"contents\":"+"\""+contentPayload+"\"")
	}
	if len(attributes) > 0 {
		a := []string{}
		for key, value := range attributes {
			a = append(a, "\""+key+"\":\""+value+"\"")
		}
		attributesStr := "\"attributes\":{" + strings.Join(a, ",") + "}"
		spec = append(spec, attributesStr)
	}

	return respBody, err
}

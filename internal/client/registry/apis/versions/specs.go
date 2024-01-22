// Copyright 2023 Google LLC
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

package versions

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"internal/apiclient"
)

// CreateSpec
func CreateSpec(apiName string, apiVersion string, apiSpecId string, name string, fileName string,
	description string, sourceURI string, contents string, labels map[string]string,
	annotations map[string]string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs")
	q := u.Query()
	q.Set("apiSpecId", apiSpecId)
	u.RawQuery = q.Encode()

	specContent := []string{}
	specContent = append(specContent, "\"name\":"+"\""+name+"\"")
	if description != "" {
		specContent = append(specContent, "\"description\":"+"\""+description+"\"")
	}
	if fileName != "" {
		specContent = append(specContent, "\"filename\":"+"\""+fileName+"\"")
	}
	if sourceURI != "" {
		specContent = append(specContent, "\"sourceUri\":"+"\""+sourceURI+"\"")
	}
	if contents != "" {
		specContent = append(specContent, "\"contents\":"+"\""+contents+"\"")
	}
	if len(labels) > 0 {
		l := []string{}
		for key, value := range labels {
			l = append(l, "\""+key+"\":\""+value+"\"")
		}
		labelStr := "\"labels\":{" + strings.Join(l, ",") + "}"
		specContent = append(specContent, labelStr)
	}

	if len(annotations) > 0 {
		a := []string{}
		for key, value := range annotations {
			a = append(a, "\""+key+"\":\""+value+"\"")
		}
		annotationStr := "\"annotations\":{" + strings.Join(a, ",") + "}"
		specContent = append(specContent, annotationStr)
	}
	payload := "{" + strings.Join(specContent, ",") + "}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// TagSpec
func TagSpec(apiName string, apiVersion string, name string, tag string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs", name+":tagRevision")
	payload := "{\"tag\":\"" + tag + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// RollbackSpec
func RollbackSpec(apiName string, apiVersion string, name string, revisionID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs", name+":rollback")
	payload := "{\"revisionId\":\"" + revisionID + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// DeleteSpec
func DeleteSpec(apiName string, apiVersion string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// DeleteSpecRevision
func DeleteSpecRevision(apiName string, apiVersion string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs", name+":deleteRevision")
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// GetSpec
func GetSpec(apiName string, apiVersion string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetSpecContents
func GetSpecContents(apiName string, apiVersion string, name string) (err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	var specMap map[string]interface{}
	var specFileName string
	respBody, err := GetSpec(apiName, apiVersion, name)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &specMap)
	if err != nil {
		return err
	}
	if specMap["filename"] != "" {
		specFileName = fmt.Sprintf("%s", specMap["filename"])
	} else {
		specFileName = name + ".txt"
	}
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs", name+":getContents")
	return apiclient.DownloadResource(u.String(), specFileName, "", true)
}

// ListSpecs
func ListSpecs(apiName string, apiVersion string, pageSize int,
	pageToken string, filter string, orderBy string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs")
	q := u.Query()
	if pageSize != -1 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	if orderBy != "" {
		q.Set("orderBy", orderBy)
	}

	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

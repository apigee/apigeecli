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

package env

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

// generateUploadURL
func generateUploadURL() (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "archiveDeployments:generateUploadUrl")
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

// CreatetArchive
func CreateArchive(name string, zipfile string) (respBody []byte, err error) {
	genUrlJson := make(map[string]interface{})

	genUrlResp, err := generateUploadURL()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(genUrlResp, &genUrlJson)
	if err != nil {
		return nil, err
	}

	gcsURI := fmt.Sprintf("%v", genUrlJson["uploadUri"])
	gcsurl, err := url.Parse(gcsURI)
	if err != nil {
		return nil, err
	}

	gcs_uri := fmt.Sprintf("%s://%s%s", gcsurl.Scheme, gcsurl.Host, gcsurl.Path)

	headers := make(map[string]string)
	headers["content-type"] = "application/zip"
	headers["x-goog-content-length-range"] = "0,1073741824"

	err = apiclient.PostHttpZip(false, "PUT", gcsURI, headers, zipfile)
	if err != nil {
		return nil, err
	}

	archive := []string{}
	archive = append(archive, "\"name\":\""+name+"\"")
	archive = append(archive, "\"gcs_uri\":\""+gcs_uri+"\"")

	payload := "{" + strings.Join(archive, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "archiveDeployments")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// GetArchive
func GetArchive(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "archiveDeployments", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListArchives
func ListArchives() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "archiveDeployments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// DeleteArchive
func DeleteArchive(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "archiveDeployments", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

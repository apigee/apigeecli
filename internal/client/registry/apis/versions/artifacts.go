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
	"net/url"
	"path"

	"internal/apiclient"
	"internal/client/registry/common"
)

type Action uint8

const (
	CREATE Action = iota
	UPDATE
)

// CreateArtifact
func CreateArtifact(apiName string, apiVersion string, artifactID string, name string, contents string,
	labels map[string]string, annotations map[string]string,
) (respBody []byte, err error) {
	return createOrReplaceArtifact(apiName, apiVersion, artifactID, name, contents, labels, annotations, CREATE)
}

// ReaplceArtifact
func ReplaceArtifact(apiName string, apiVersion string, artifactID string, name string, contents string,
	labels map[string]string, annotations map[string]string,
) (respBody []byte, err error) {
	return createOrReplaceArtifact(apiName, apiVersion, artifactID, name, contents, labels, annotations, UPDATE)
}

// DeleteArtifact
func DeleteArtifact(apiName string, apiVersion string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "artifacts", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// GetArtifact
func GetArtifact(apiName string, apiVersion string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "artifacts", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetArtifactContents
func GetArtifactContents(apiName string, apiVersion string, name string) (err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "artifacts", name+":getContents")
	return apiclient.DownloadResource(u.String(), name+".txt", "", true)
}

// ListArtifacts
func ListArtifacts(apiName string, apiVersion string, pageSize int,
	pageToken string, filter string, orderBy string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "artifacts")
	return common.ListArtifacts(u, pageSize, pageToken, filter, orderBy)
}

func createOrReplaceArtifact(apiName, apiVersion, artifactID string, name string, contents string,
	labels map[string]string, annotations map[string]string, action Action,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	payload := common.GetArtifactPayload(name, contents, labels, annotations)

	if action == CREATE {
		u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "artifacts")
		q := u.Query()
		q.Set("artifactId", artifactID)
		u.RawQuery = q.Encode()
		respBody, err = apiclient.HttpClient(u.String(), payload)
	} else {
		u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "artifacts", artifactID)
		respBody, err = apiclient.HttpClient(u.String(), payload, "PUT")
	}
	return respBody, err
}

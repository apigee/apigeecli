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

// CreateSpecArtifact
func CreateSpecArtifact(apiName string, apiVersion string, specName string, artifactId string,
	name string, contents string, labels map[string]string,
	annotations map[string]string,
) (respBody []byte, err error) {
	return createOrReplaceSpecArtifact(apiName, apiVersion, specName, artifactId, name, contents, labels, annotations, CREATE)
}

// ReplaceSpecArtifact
func ReplaceSpecArtifact(apiName string, apiVersion string, specName string, artifactId string,
	name string, contents string, labels map[string]string,
	annotations map[string]string,
) (respBody []byte, err error) {
	return createOrReplaceSpecArtifact(apiName, apiVersion, specName, artifactId, name, contents, labels, annotations, UPDATE)
}

// DeleteSpecArtifact
func DeleteSpecArtifact(apiName string, apiVersion string, specName string,
	name string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs",
		specName, "artifacts", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// GetSpecArtifact
func GetSpecArtifact(apiName string, apiVersion string, specName string,
	name string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs",
		specName, "artifacts", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetSpecArtifactContents
func GetSpecArtifactContents(apiName string, apiVersion string, specName string,
	name string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion,
		"specs", specName, "artifacts", name+":getContents")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListSpecArtifacts
func ListSpecArtifacts(apiName string, apiVersion string, specName string, pageSize int,
	pageToken string, filter string, orderBy string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs",
		specName, "artifacts")
	return common.ListArtifacts(u, pageSize, pageToken, filter, orderBy)
}

func createOrReplaceSpecArtifact(apiName string, apiVersion string, specName string, artifactID string,
	name string, contents string, labels map[string]string,
	annotations map[string]string, action Action,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	payload := common.GetArtifactPayload(name, contents, labels, annotations)
	if action == CREATE {
		u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs", specName, "artifacts")
		q := u.Query()
		q.Set("artifactId", artifactID)
		u.RawQuery = q.Encode()
		respBody, err = apiclient.HttpClient(u.String(), payload)
	} else {
		u.Path = path.Join(u.Path, "apis", apiName, "versions", apiVersion, "specs", specName, "artifacts", artifactID)
		respBody, err = apiclient.HttpClient(u.String(), payload, "PUT")
	}
	return respBody, err
}

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

package apis

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"internal/apiclient"
)

// CreateDeployment
func CreateDeployment(apiDeploymentID string, name string, displayName string, description string,
	apiSpecRevision string, endpointURI string, externalChannelURI string, intendedAudience string,
	accessGuidance string, labels map[string]string, annotations map[string]string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiDeploymentID, "deployments")
	q := u.Query()
	q.Set("apiDeploymentId", name)
	u.RawQuery = q.Encode()
	apiDeployment := []string{}
	apiDeployment = append(apiDeployment, "\"name\":"+"\""+name+"\"")
	if displayName != "" {
		apiDeployment = append(apiDeployment, "\"displayName\":"+"\""+displayName+"\"")
	}
	if description != "" {
		apiDeployment = append(apiDeployment, "\"description\":"+"\""+description+"\"")
	}
	if apiSpecRevision != "" {
		apiDeployment = append(apiDeployment, "\"apiSpecRevision\":"+"\""+apiSpecRevision+"\"")
	}
	if endpointURI != "" {
		apiDeployment = append(apiDeployment, "\"endpointUri\":"+"\""+endpointURI+"\"")
	}
	if externalChannelURI != "" {
		apiDeployment = append(apiDeployment, "\"externalChannelUri\":"+"\""+externalChannelURI+"\"")
	}
	if intendedAudience != "" {
		apiDeployment = append(apiDeployment, "\"intendedAudience\":"+"\""+intendedAudience+"\"")
	}
	if accessGuidance != "" {
		apiDeployment = append(apiDeployment, "\"accessGuidance\":"+"\""+accessGuidance+"\"")
	}
	if len(labels) > 0 {
		l := []string{}
		for key, value := range labels {
			l = append(l, "\""+key+"\":\""+value+"\"")
		}
		labelStr := "\"labels\":{" + strings.Join(l, ",") + "}"
		apiDeployment = append(apiDeployment, labelStr)
	}

	if len(annotations) > 0 {
		a := []string{}
		for key, value := range annotations {
			a = append(a, "\""+key+"\":\""+value+"\"")
		}
		annotationStr := "\"annotations\":{" + strings.Join(a, ",") + "}"
		apiDeployment = append(apiDeployment, annotationStr)
	}
	payload := "{" + strings.Join(apiDeployment, ",") + "}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Tag
func Tag(name string, deployment string, tag string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "deployments", deployment+":tagRevision")
	payload := "{\"tag\":\"" + tag + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Rollback
func Rollback(name string, deployment string, revisionID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "deployments", deployment+":rollback")
	payload := "{\"revisionId\":\"" + revisionID + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// DeleteDeployment
func DeleteDeployment(name string, deployment string, force bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "deployments", deployment)
	if force {
		q := u.Query()
		q.Set("force", strconv.FormatBool(force))
		u.RawQuery = q.Encode()
	}
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// DeleteDeploymentRevision
func DeleteDeploymentRevision(name string, deployment string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "deployments", deployment+":deleteRevision")
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// GetDeployment
func GetDeployment(name string, deployment string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "deployments", deployment)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListDeploymentRevisions
func ListDeploymentRevisions(apiName string, name string, pageSize int, pageToken string,
	filter string, orderBy string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiName, "deployments", name+":listRevisions")
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

// ListDeployments
func ListDeployments(name string, pageSize int, pageToken string,
	filter string, orderBy string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", name, "deployments")
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

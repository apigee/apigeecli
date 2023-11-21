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

package datastores

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

type axdatastores struct {
	Datastores []axdatastore `json:"datastores,omitempty"`
}

type axdatastore struct {
	Self            string          `json:"self,omitempty"`
	DisplayName     string          `json:"displayName,omitempty"`
	Org             string          `json:"org,omitempty"`
	TargetType      string          `json:"targetType,omitempty"`
	CreateTime      string          `json:"createTime,omitempty"`
	LastUpdateTime  string          `json:"lastUpdateTime,omitempty"`
	DatastoreConfig datastoreconfig `json:"datastoreConfig,omitempty"`
}

type datastoreconfig struct {
	ProjectID   string `json:"projectId,omitempty"`
	BucketName  string `json:"bucketName,omitempty"`
	Path        string `json:"path,omitempty"`
	TablePrefix string `json:"tablePrefix,omitempty"`
	DatasetName string `json:"datasetName,omitempty"`
}

type Action uint8

const (
	CREATE Action = iota
	UPDATE
	TEST
)

// Create
func Create(displayName string, targetType string, projectID string,
	bucketName string, gcsPath string, datasetName string,
	tablePrefix string,
) (respBody []byte, err error) {
	return createOrTestorUpdate("", displayName, targetType, projectID, bucketName,
		gcsPath, datasetName, tablePrefix, CREATE)
}

// Update
func Update(id string, displayName string, targetType string, projectID string,
	bucketName string, gcsPath string, datasetName string,
	tablePrefix string,
) (respBody []byte, err error) {
	return createOrTestorUpdate(id, displayName, targetType, projectID, bucketName,
		gcsPath, datasetName, tablePrefix, UPDATE)
}

// Test
func Test(id string, displayName string, targetType string, projectID string,
	bucketName string, gcsPath string, datasetName string,
	tablePrefix string,
) (respBody []byte, err error) {
	return createOrTestorUpdate(id, displayName, targetType, projectID, bucketName,
		gcsPath, datasetName, tablePrefix, TEST)
}

func createOrTestorUpdate(id string, displayName string, targetType string, projectID string,
	bucketName string, gcsPath string, datasetName string,
	tablePrefix string, action Action,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	ds := []string{}

	ds = append(ds, "\"displayName\":"+"\""+displayName+"\"")
	ds = append(ds, "\"targetType\":"+"\""+targetType+"\"")

	dsc := []string{}
	dsc = append(dsc, "\"projectId\":"+"\""+projectID+"\"")
	if targetType == "gcs" {
		if bucketName == "" || gcsPath == "" {
			return nil, fmt.Errorf("bucketName and path are mandatory parameters")
		}
		dsc = append(dsc, "\"bucketName\":"+"\""+bucketName+"\"")
		dsc = append(dsc, "\"path\":"+"\""+gcsPath+"\"")
	} else if targetType == "bigquery" {
		if datasetName == "" || tablePrefix == "" {
			return nil, fmt.Errorf("datasetName and tablePrefix are mandatory parameters")
		}
		dsc = append(dsc, "\"datasetName\":"+"\""+datasetName+"\"")
		dsc = append(dsc, "\"tablePrefix\":"+"\""+tablePrefix+"\"")
	}

	dataConfig := "{" + strings.Join(dsc, ",") + "}"
	ds = append(ds, "\"datastoreConfig\":"+dataConfig)

	payload := "{" + strings.Join(ds, ",") + "}"

	if action == CREATE {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "analytics", "datastores")
		respBody, err = apiclient.HttpClient(u.String(), payload)
	} else if action == UPDATE {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "analytics", "datastores", id)
		respBody, err = apiclient.HttpClient(u.String(), payload, "PUT")
	} else if action == TEST {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "analytics", "datastores:test")
		respBody, err = apiclient.HttpClient(u.String(), payload)
	}

	return respBody, err
}

// Delete
func Delete(id string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "analytics", "datastores", id)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get(id string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "analytics", "datastores", id)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func GetName(name string) (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	listBody, err := List("")
	if err != nil {
		return nil, err
	}

	listDatastores := axdatastores{}

	err = json.Unmarshal(listBody, &listDatastores)
	if err != nil {
		return nil, err
	}

	if len(listDatastores.Datastores) < 1 {
		return nil, fmt.Errorf("data store was not found")
	}

	for _, ds := range listDatastores.Datastores {
		if ds.DisplayName == name {
			return json.Marshal(ds)
		}
	}
	return nil, fmt.Errorf("data store was not found")
}

// GetVersion
func GetVersion(name string) (version string, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	listBody, err := List("")
	if err != nil {
		return "", err
	}

	listDatastores := axdatastores{}

	err = json.Unmarshal(listBody, &listDatastores)
	if err != nil {
		return "", err
	}

	if len(listDatastores.Datastores) < 1 {
		return "", fmt.Errorf("data store was not found")
	}

	for _, ds := range listDatastores.Datastores {
		if ds.DisplayName == name {
			return getVersion(ds.Self), nil
		}
	}
	return "", fmt.Errorf("data store was not found")
}

// List
func List(targetType string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "analytics", "datastores")
	if (targetType != "") && (targetType != "gcs" || targetType != "bigquery") {
		return nil, fmt.Errorf("invalid targetType. Must be gcs or bigquery")
	}
	q := u.Query()
	if targetType != "" {
		q.Set("targetType", targetType)
	}

	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// getVersion
func getVersion(name string) (version string) {
	s := strings.Split(name, "/")
	return s[len(s)-1]
}

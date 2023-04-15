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

package datacollectors

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"internal/apiclient"
)

type dcollectors struct {
	DataCollector []datacollector `json:"dataCollectors,omitempty"`
}

type datacollector struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Type           string `json:"type,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`
}

// Create
func Create(name string, description string, collectorType string) (respBody []byte, err error) {
	datacollector := []string{}

	datacollector = append(datacollector, "\"name\":\""+name+"\"")
	if description != "" {
		datacollector = append(datacollector, "\"description\":\""+description+"\"")
	}
	datacollector = append(datacollector, "\"type\":\""+collectorType+"\"")

	payload := "{" + strings.Join(datacollector, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "datacollectors")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "datacollectors", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "datacollectors", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "datacollectors")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Import
func Import(filePath string) (err error) {
	var dCollectors dcollectors

	if dCollectors, err = readDataCollectorsFile(filePath); err != nil {
		return err
	}

	if len(dCollectors.DataCollector) < 1 {
		return nil
	}

	for _, dCollector := range dCollectors.DataCollector {
		if _, err = Create(dCollector.Name, dCollector.Description, dCollector.Type); err != nil {
			return err
		}
	}
	return nil
}

func readDataCollectorsFile(filePath string) (dcollectors, error) {
	dCollectors := dcollectors{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return dCollectors, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return dCollectors, err
	}

	err = json.Unmarshal(byteValue, &dCollectors)

	if err != nil {
		return dCollectors, err
	}

	return dCollectors, nil
}

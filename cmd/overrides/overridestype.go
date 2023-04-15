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

package overrides

import (
	"encoding/json"
	"os"

	"internal/clilog"

	yaml "gopkg.in/yaml.v2"
)

type overrides struct {
	Org          string        `yaml:"org,omitempty"`
	Gcp          gcp           `yaml:"gcp,omitempty"`
	Envs         []environment `yaml:"envs,omitempty"`
	Virtualhosts []virtualhost `yaml:"virtualhosts,omitempty"`
}

type gcp struct {
	ProjectID string `yaml:"projectID,omitempty"`
	Region    string `yaml:"region,omitempty"`
}

type environment struct {
	Name                string               `yaml:"name,omitempty"`
	ServiceAccountPaths *serviceAccountPaths `yaml:"serviceAccountPaths,omitempty"`
}

type serviceAccountPaths struct {
	Synchronizer string `yaml:"synchronizer,omitempty"`
}

type virtualhost struct {
	Name string `yaml:"name,omitempty"`
}

var hybridOverrides = overrides{}

func readOverrides(filepath string) (err error) {
	source, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(source, &hybridOverrides)
	if err != nil {
		return err
	}
	return
}

func getOrg() string {
	return hybridOverrides.Org
}

func getOrgRegion() string {
	return hybridOverrides.Gcp.Region
}

func getEnvs() []string {
	environmentList := []string{}
	for _, environment := range hybridOverrides.Envs {
		environmentList = append(environmentList, environment.Name)
	}
	return environmentList
}

func getEnvGroups() []string {
	environmentGroupList := []string{}
	for _, environmentGroup := range hybridOverrides.Virtualhosts {
		environmentGroupList = append(environmentGroupList, environmentGroup.Name)
	}
	return environmentGroupList
}

func getSyncServiceAccounts() []string {
	syncIdentityList := []string{}
	for _, environment := range hybridOverrides.Envs {
		if environment.ServiceAccountPaths != nil && environment.ServiceAccountPaths.Synchronizer != "" {
			identity, _ := getIdentity(environment.ServiceAccountPaths.Synchronizer)
			syncIdentityList = append(syncIdentityList, identity)
		}
	}
	return syncIdentityList
}

func getIdentity(serviceAccountPath string) (string, error) {
	serviceAccount := make(map[string]interface{})
	saBytes, err := os.ReadFile(serviceAccountPath)
	if err != nil {
		clilog.Error.Println(err)
		return "", err
	}
	if err := json.Unmarshal(saBytes, &serviceAccount); err != nil {
		clilog.Error.Println(err)
		return "", err
	}
	return serviceAccount["client_email"].(string), nil
}

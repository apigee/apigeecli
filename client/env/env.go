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

package env

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/srinandan/apigeecli/apiclient"
)

//Get
func Get(config bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if config {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "deployedConfig")
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	}
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//SetEnvProperty is used to set env properties
func SetEnvProperty(name string, value string) (err error) {
	//EnvProperty contains an individual org flag or property
	type envProperty struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}
	//EnvProperties stores all the org feature flags and properties
	type envProperties struct {
		Property []envProperty `json:"property,omitempty"`
	}

	//Env structure
	type environment struct {
		Name           string        `json:"name,omitempty"`
		Description    string        `json:"description,omitempty"`
		CreatedAt      string        `json:"-,omitempty"`
		LastModifiedAt string        `json:"-,omitempty"`
		Properties     envProperties `json:"properties,omitempty"`
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	//get env details
	envBody, err := apiclient.HttpClient(false, u.String())
	if err != nil {
		return err
	}

	env := environment{}
	err = json.Unmarshal(envBody, &env)
	if err != nil {
		return err
	}

	//check if the property exists
	found := false
	for i, properties := range env.Properties.Property {
		if properties.Name == name {
			fmt.Println("Property found, enabling property")
			env.Properties.Property[i].Value = value
			found = true
			break
		}
	}

	if !found {
		//set the property
		newProp := envProperty{}
		newProp.Name = name
		newProp.Value = value

		env.Properties.Property = append(env.Properties.Property, newProp)
	}

	newEnvBody, err := json.Marshal(env)
	if err != nil {
		return err
	}

	u, _ = url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	_, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(newEnvBody), "PUT")

	return err
}

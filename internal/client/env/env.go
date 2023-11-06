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
	"strings"

	"internal/apiclient"
	"internal/clilog"
)

// Create
func Create(deploymentType string, fwdProxyURI string) (respBody []byte, err error) {
	environment := []string{}
	environment = append(environment, "\"name\":\""+apiclient.GetApigeeEnv()+"\"")

	if deploymentType != "" {
		if deploymentType != "PROXY" && deploymentType != "ARCHIVE" {
			return nil, fmt.Errorf("deploymentType must be PROXY or ARCHIVE")
		}
		environment = append(environment, "\"deployment_type\":\""+deploymentType+"\"")
	}

	if fwdProxyURI != "" {
		environment = append(environment, "\"forwardProxyUri\":\""+fwdProxyURI+"\"")
	}

	payload := "{" + strings.Join(environment, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Delete
func Delete() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get(config bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if config {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "deployedConfig")
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	}
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetDeployments
func GetDeployments(sharedflows bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if sharedflows {
		q := u.Query()
		q.Set("sharedFlows", "true")
		u.RawQuery = q.Encode()
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "deployments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func GetAllDeployments() (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	proxiesResponse, err := GetDeployments(false)
	if err != nil {
		return nil, err
	}

	sharedFlowsResponse, err := GetDeployments(true)
	if err != nil {
		return nil, err
	}

	deployments := []string{}

	deployments = append(deployments, "\"proxies\":"+string(proxiesResponse))
	deployments = append(deployments, "\"sharedFlows\":"+string(sharedFlowsResponse))
	payload := "{" + strings.Join(deployments, ",") + "}"

	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	err = apiclient.PrettyPrint("json", []byte(payload))
	return []byte(payload), err
}

// GetDeployedConfig
func GetDeployedConfig() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "deployedConfig")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// SetEnvProperty is used to set env properties
func SetEnvProperty(name string, value string) (err error) {
	// EnvProperty contains an individual org flag or property
	type envProperty struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}
	// EnvProperties stores all the org feature flags and properties
	type envProperties struct {
		Property []envProperty `json:"property,omitempty"`
	}

	// Env structure
	type environment struct {
		Name           string        `json:"name,omitempty"`
		Description    string        `json:"description,omitempty"`
		CreatedAt      string        `json:"-,omitempty"`
		LastModifiedAt string        `json:"-,omitempty"`
		Properties     envProperties `json:"properties,omitempty"`
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	// get env details
	apiclient.ClientPrintHttpResponse.Set(false)
	envBody, err := apiclient.HttpClient(u.String())
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	if err != nil {
		return err
	}

	env := environment{}
	err = json.Unmarshal(envBody, &env)
	if err != nil {
		return err
	}

	// check if the property exists
	found := false
	for i, properties := range env.Properties.Property {
		if properties.Name == name {
			clilog.Info.Println("Property found, enabling property")
			env.Properties.Property[i].Value = value
			found = true
			break
		}
	}

	if !found {
		// set the property
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
	_, err = apiclient.HttpClient(u.String(), string(newEnvBody), "PUT")

	return err
}

// ClearEnvProperties is used to set env properties
func ClearEnvProperties() (err error) {
	// EnvProperty contains an individual org flag or property
	type envProperty struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}
	// EnvProperties stores all the org feature flags and properties
	type envProperties struct {
		Property []envProperty `json:"property,omitempty"`
	}

	// Env structure
	type environment struct {
		Name           string        `json:"name,omitempty"`
		Description    string        `json:"description,omitempty"`
		CreatedAt      string        `json:"-,omitempty"`
		LastModifiedAt string        `json:"-,omitempty"`
		Properties     envProperties `json:"-,omitempty"`
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	// get env details
	apiclient.ClientPrintHttpResponse.Set(false)
	envBody, err := apiclient.HttpClient(u.String())
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	if err != nil {
		return err
	}

	env := environment{}
	err = json.Unmarshal(envBody, &env)
	if err != nil {
		return err
	}

	newEnv := environment{}
	newEnv.Name = env.Name
	newEnv.Description = env.Description
	newEnvBody, err := json.Marshal(newEnv)
	if err != nil {
		return err
	}

	u, _ = url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	_, err = apiclient.HttpClient(u.String(), string(newEnvBody), "PUT")

	return err
}

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
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"internal/apiclient"
	"internal/clilog"
)

type Environments struct {
	Environment []Environment `json:"environment,omitempty"`
}

type Environment struct {
	Name            string     `json:"name,omitempty"`
	DisplayName     string     `json:"displayName,omitempty"`
	DeploymentType  string     `json:"deploymentType,omitempty"`
	ApiProxyType    string     `json:"apiProxyType,omitempty"`
	ForwardProxyUri string     `json:"forwardProxyUri,omitempty"`
	Type            string     `json:"type,omitempty"`
	Properties      []Property `json:"properties,omitempty"`
}

type Property struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type AddonsConfig struct {
	AnalyticsConfig AnalyticsConfig `json:"analyticsConfig,omitempty"`
}

type AnalyticsConfig struct {
	Enabled         bool   `json:"enabled,omitempty"`
	ExpireTimeMills string `json:"expireTimeMills,omitempty"`
	State           string `json:"state,omitempty"`
	UpdateTime      string `json:"updateTime,omitempty"`
}

// Create
func Create(envType string, deploymentType string, fwdProxyURI string) (respBody []byte, err error) {
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

	if envType != "" {
		environment = append(environment, "\"type\":\""+envType+"\"")
	}

	payload := "{" + strings.Join(environment, ",") + "}"

	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Delete
func Delete() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get(config bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
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
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetDeployments
func GetDeployments(sharedflows bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
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
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
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

	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
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

	u, _ = url.Parse(apiclient.GetApigeeBaseURL())
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

	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
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

	u, _ = url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv())
	_, err = apiclient.HttpClient(u.String(), string(newEnvBody), "PUT")

	return err
}

// GetSecurityActionsConfig
func GetSecurityActionsConfig() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "securityActionsConfig")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetSecurityRuntimeConfig
func GetSecurityRuntimeConfig() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "apiSecurityRuntimeConfig")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// UpdateSecurityActionsConfig
func UpdateSecurityActionsConfig(action bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "apiSecurityRuntimeConfig")

	payload := fmt.Sprintf("{ \"name\": \"organizations/%s/environments/%s/securityActionsConfig\",\"payload\": %t}",
		apiclient.GetApigeeOrg(), apiclient.GetApigeeEnv(), action)

	respBody, err = apiclient.HttpClient(u.String(), payload, "PATCH")
	return respBody, err
}

// Export
func Export() (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	var envList []string
	environmentList := Environments{}

	envRespBody, err := List()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(envRespBody, &envList)
	if err != nil {
		return nil, err
	}
	for _, e := range envList {
		apiclient.SetApigeeEnv(e)
		envRespBody, err := Get(false)
		if err != nil {
			return nil, err
		}
		environ := Environment{}
		err = json.Unmarshal(envRespBody, &environ)
		if err != nil {
			return nil, err
		}
		environmentList.Environment = append(environmentList.Environment, environ)
	}
	respBody, err = json.Marshal(&environmentList)
	if err != nil {
		return nil, err
	}
	respBody, err = apiclient.PrettifyJSON(respBody)
	return respBody, err
}

func Import(filePath string) (err error) {
	entities, err := readEnvironmentsFile(filePath)
	var errs []error
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	numEntities := len(entities.Environment)
	clilog.Debug.Printf("Found %d environments in the file\n", numEntities)

	for _, entity := range entities.Environment {
		clilog.Info.Printf("Creating environment %s\n", entity.Name)
		apiclient.SetApigeeEnv(entity.Name)
		_, err = Create(entity.Type, entity.DeploymentType, entity.ForwardProxyUri)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func readEnvironmentsFile(filePath string) (environmentList Environments, err error) {
	environmentList = Environments{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return environmentList, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return environmentList, err
	}

	err = json.Unmarshal(byteValue, &environmentList)
	if err != nil {
		return environmentList, err
	}

	return environmentList, nil
}

func MarshalEnvironmentList(contents []byte) (envronmentList Environments, err error) {
	if err = json.Unmarshal(contents, &envronmentList); err != nil {
		return envronmentList, err
	}
	return envronmentList, nil
}

func GetAddonsConfig(environment string) (respBody []byte, err error) {

	if environment == "" {
		environment = apiclient.GetApigeeEnv()
	}

	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		environment, "addonsConfig")

	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func IsAnalyticsEnabled(environment string) (enabled bool, err error) {
	var tmp string
	apiclient.ClientPrintHttpResponse.Set(false)
	tmp = apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv(environment)
	respBody, err := Get(false)
	apiclient.SetApigeeEnv(tmp)
	if err != nil {
		return false, err
	}
	e := Environment{}

	if err = json.Unmarshal(respBody, &e); err != nil {
		return false, err
	}

	if e.Type == "BASE" {
		return false, nil
	}

	respBody, err = GetAddonsConfig(environment)
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	if err != nil {
		return false, err
	}
	addonsConfig := AddonsConfig{}
	err = json.Unmarshal(respBody, &addonsConfig)
	if err != nil {
		return false, err
	}
	return addonsConfig.AnalyticsConfig.Enabled, nil
}

// Copyright 2024 Google LLC
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

package hub

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"internal/apiclient"
	"internal/clilog"
	"internal/cmd/utils"
)

type DeploymentType string

const (
	APIGEE              DeploymentType = "apigee"
	APIGEE_HYBRID       DeploymentType = "apigee-hybrid"
	APIGEE_EDGE_PRIVATE DeploymentType = "apigee-edge-private"
	APIGEE_EDGE_PUBLIC  DeploymentType = "apigee-edge-public"
	MOCK_SERVER         DeploymentType = "mock-server"
	CLOUD_API_GATEWAY   DeploymentType = "cloud-api-gateway"
	CLOUD_ENDPOINTS     DeploymentType = "cloud-endpoints"
	UNMANAGED           DeploymentType = "unmanaged"
	OTHERS              DeploymentType = "others"
)

type spec struct {
	DisplayName   string             `json:"displayName,omitempty"`
	SpecType      enumAttributeValue `json:"specType,omitempty"`
	SourceURI     string             `json:"sourceUri,omitempty"`
	Contents      content            `json:"contents,omitempty"`
	Documentation documentation      `json:"documentation,omitempty"`
	ParsingMode   string             `json:"parsingMode,omitempty"`
}

type EnvironmentType string

const (
	DEVELOPMENT   EnvironmentType = "development"
	STAGING       EnvironmentType = "staging"
	TEST          EnvironmentType = "test"
	PREPRODUCTION EnvironmentType = "pre-prod"
	PRODUCTION    EnvironmentType = "prod"
)

type SloType string

const (
	SLO99_99 SloType = "99-99"
	SLO99_95 SloType = "99-95"
	SLO99_9  SloType = "99-9"
	SLO99_5  SloType = "99-5"
)

type documentation struct {
	ExternalUri string `json:"externalUri,omitempty"`
}

type content struct {
	Contents string `json:"contents,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

type allowedValue struct {
	Id          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
	Immutable   bool   `json:"immutable"`
}

type enumValues struct {
	Values []allowedValue `json:"values,omitempty"`
}

type enumAttributeValue struct {
	EnumValues enumValues `json:"enumValues,omitempty"`
}

type Action uint8

const (
	CREATE Action = iota
	UPDATE
)

func CreateInstance(apiHubInstanceId string, cmekName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apiHubInstances")

	if apiHubInstanceId != "" {
		q := u.Query()
		q.Set("apiHubInstanceId", apiHubInstanceId)
		u.RawQuery = q.Encode()
	}

	instance := []string{}
	config := []string{}

	name := fmt.Sprintf("projects/%s/locations/%s/instance", apiclient.GetProjectID(), apiclient.GetRegistryRegion())
	instance = append(instance, "\"name\":"+"\""+name+"\"")

	config = append(config, "\"cmekKeyName\":"+"\""+cmekName+"\"")
	configJson := "{" + strings.Join(config, ",") + "}"
	instance = append(instance, "\"config\":"+configJson)

	payload := "{" + strings.Join(instance, ",") + "}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func GetInstance(apiHubInstanceId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apiHubInstances", apiHubInstanceId)

	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func LookupInstance() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apiHubInstances:lookup")

	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func RegisterHostProject(registrationId string, gcpProjectId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "hostProjectRegistrations")

	q := u.Query()
	q.Set("hostProjectRegistrationId", registrationId)
	u.RawQuery = q.Encode()

	payload := fmt.Sprintf("{\"gcpProject\":\"%s\"}", gcpProjectId)
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func ListHostProjects(filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list("hostProjectRegistrations", filter, pageSize, pageToken)
}

func CreateRuntimeProjectAttachment(runtimeProjectAttachmentId string, runtimeProject string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "runtimeProjectAttachments")

	q := u.Query()
	q.Set("runtimeProjectAttachmentId", runtimeProjectAttachmentId)
	u.RawQuery = q.Encode()

	payload := fmt.Sprintf("{\"runtimeProject\":\"%s\"}", runtimeProject)
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func DeleteRuntimeProjectAttachment(runtimeProjectAttachmentId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "runtimeProjectAttachments", runtimeProjectAttachmentId)

	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func GetRuntimeProjectAttachment(runtimeProjectAttachmentId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "runtimeProjectAttachments", runtimeProjectAttachmentId)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func ListRuntimeProjectAttachments(filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list("runtimeProjectAttachments", filter, pageSize, pageToken)
}

func CreateApi(apiID string, contents []byte) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis")

	q := u.Query()
	q.Set("apiId", apiID)
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), string(contents))
	return respBody, err
}

func DeleteApi(apiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID)

	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func GetApi(apiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID)

	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func ListApi(filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list("apis", filter, pageSize, pageToken)
}

func UpdateApi(apiID string, contents []byte) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID)
	q := u.Query()
	q.Set("updateMask",
		"display_name,description,owner,documentation,target_user,team,business_unit,maturity_level,attributes")
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String(), string(contents), "PATCH")
	return respBody, err
}

func ExportApi(apiID string, folder string) (err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	folderName := path.Join(folder, "api"+utils.DefaultFileSplitter+apiID)

	// create a folder for the api
	_, err = os.Stat(folderName)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(folderName, 0o755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// write api details
	apiBody, err := GetApi(apiID)
	if err != nil {
		return err
	}
	apiFileName := fmt.Sprintf("api%s%s.json", utils.DefaultFileSplitter, apiID)
	err = apiclient.WriteByteArrayToFile(path.Join(folderName, apiFileName), false, apiBody)
	if err != nil {
		return err
	}

	type apiVersions struct {
		Versions      []map[string]interface{} `json:"versions"`
		NextPageToken string                   `json:"nextPageToken"`
	}

	apiVersionsObj := apiVersions{}

	for {

		a := apiVersions{}

		apiVersionsBody, err := ListApiVersions(apiID, "", -1, "")
		if err != nil {
			return fmt.Errorf("Error listing api versions: %s", err)
		}
		err = json.Unmarshal(apiVersionsBody, &a)
		if err != nil {
			return fmt.Errorf("Error unmarshalling api versions: %s", err)
		}

		apiVersionsObj.Versions = append(apiVersionsObj.Versions, a.Versions...)
		if a.NextPageToken == "" {
			break
		}
	}

	for _, version := range apiVersionsObj.Versions {
		versionName := filepath.Base(version["name"].(string))
		apiVersionBody, err := GetApiVersion(versionName, apiID)
		if err != nil {
			return fmt.Errorf("Error getting api version: %s", err)
		}
		apiVersionFileName := fmt.Sprintf("api%s%s%s%s.json",
			utils.DefaultFileSplitter, apiID, utils.DefaultFileSplitter, versionName)
		err = apiclient.WriteByteArrayToFile(path.Join(folderName, apiVersionFileName), false, apiVersionBody)
		if err != nil {
			return err
		}

		if err = ExportApiVersionSpecs(apiID, versionName, folderName); err != nil {
			return err
		}
	}

	return nil
}

func CreateApiVersion(versionID string, apiID string, contents []byte) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions")
	q := u.Query()
	if versionID != "" {
		q.Set("versionId", versionID)
	}
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String(), string(contents))
	return respBody, err
}

func GetApiVersion(versionID, apiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteApiVersion(versionID, apiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func ListApiVersions(apiID string, filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list(path.Join("apis", apiID, "versions"), filter, pageSize, pageToken)
}

func GetApiVersionsDefinitions(apiID string, versionID string, definition string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "definitions", definition)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func CreateApiVersionsSpec(apiID string, versionID string, specID string, displayName string,
	contents []byte, mimeType string, sourceURI string, documentation string,
) (respBody []byte, err error) {
	return createOrUpdateApiVersionSpec(apiID, versionID, specID, displayName, contents, mimeType, sourceURI, documentation, CREATE)
}

func GetApiVersionSpec(apiID string, versionID string, specID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs", specID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteApiVersionSpec(apiID string, versionID string, specID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs", specID)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func GetApiVersionsSpecContents(apiID string, versionID string, specID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs", specID, ":contents")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func ListApiVersionSpecs(apiID string, versionID string, filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list(path.Join("apis", apiID, "versions", versionID, "specs"), filter, pageSize, pageToken)
}

func LintApiVersionSpec(apiID string, versionID string, specID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs", specID, ":lint")
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

func UpdateApiVersionSpec(apiID string, versionID string, specID string, displayName string,
	contents []byte, mimeType string, sourceURI string, documentation string,
) (respBody []byte, err error) {
	return createOrUpdateApiVersionSpec(apiID, versionID, specID, displayName, contents, mimeType, sourceURI, documentation, UPDATE)
}

func ExportApiVersionSpecs(apiID string, versionID string, folder string) (err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	var specResp map[string]string

	// create a folder for the specs
	_, err = os.Stat(path.Join(folder, "specs"))
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(path.Join(folder, "specs"), 0o755)
		} else {
			return err
		}
	}

	specsBody, err := ListApiVersionSpecs(apiID, versionID, "", -1, "")
	if err != nil {
		return err
	}

	specList := getSpecIDList(specsBody)
	for _, spec := range specList {
		specBody, err := GetApiVersionsSpecContents(apiID, versionID, spec)
		if err != nil {
			return fmt.Errorf("unable to complete GetApiVersionsSpecContents")
		}

		if err = json.Unmarshal(specBody, &specResp); err != nil {
			return fmt.Errorf("unable to unmarshal specBody")
		}

		specContent, err := base64.StdEncoding.DecodeString(specResp["contents"])
		if err != nil {
			return fmt.Errorf("unable to decode contents %v", err)
		}
		err = apiclient.WriteByteArrayToFile(path.Join(folder,
			"specs", spec), false, specContent)
		if err != nil {
			return err
		}
	}

	return nil
}

func createOrUpdateApiVersionSpec(apiID string, versionID string, specID string, displayName string,
	contents []byte, mimeType string, sourceURI string, documentation string, action Action,
) (respBody []byte, err error) {
	s := spec{}
	s.DisplayName = displayName
	if documentation != "" {
		s.Documentation.ExternalUri = documentation
	}

	if contents != nil {
		s.Contents.Contents = base64.StdEncoding.EncodeToString(contents)
	}

	if strings.Contains(mimeType, "yaml") || strings.Contains(mimeType, "yml") {
		s.Contents.MimeType = "application/yaml"
		s.SpecType = getSpecType("openapi")
	} else if strings.Contains(mimeType, "json") {
		s.Contents.MimeType = "application/json"
		s.SpecType = getSpecType("openapi")
	} else if strings.Contains(mimeType, "wsdl") {
		s.Contents.MimeType = "application/wsdl"
		s.SpecType = getSpecType("wsdl")
	} else if strings.Contains(mimeType, "proto") {
		s.Contents.MimeType = "application/text"
		s.SpecType = getSpecType("proto")
	} else {
		s.Contents.MimeType = "application/text"
	}

	if sourceURI != "" {
		s.SourceURI = sourceURI
	}

	payload, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())

	if action == CREATE {
		u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs")
		q := u.Query()
		q.Set("specId", specID)
		u.RawQuery = q.Encode()
		respBody, err = apiclient.HttpClient(u.String(), string(payload))
	} else {
		u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs", specID)
		q := u.Query()
		q.Set("updateMask", "display_name,source_uri,documentation,contents,spec_type")
		u.RawQuery = q.Encode()
		respBody, err = apiclient.HttpClient(u.String(), string(payload), "PATCH")
	}
	return respBody, err
}

func CreateDependency(dependencyID string, description string, consumerDisplayName string,
	consumerOperationResourceName string, consumerExternalApiResourceName string, supplierDisplayName string,
	supplierOperationResourceName string, supplierExternalApiResourceName string,
) (respBody []byte, err error) {
	type consumer struct {
		DisplayName             string  `json:"displayName,omitempty"`
		OperationResourceName   *string `json:"operationResourceName,omitempty"`
		ExternalApiResourceName *string `json:"externalApiResourceName,omitempty"`
	}

	type supplier struct {
		DisplayName             string  `json:"displayName,omitempty"`
		OperationResourceName   *string `json:"operationResourceName,omitempty"`
		ExternalApiResourceName *string `json:"externalApiResourceName,omitempty"`
	}

	type dependency struct {
		Description string   `json:"description,omitempty"`
		Consumer    consumer `json:"consumer,omitempty"`
		Supplier    supplier `json:"supplier,omitempty"`
	}

	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "dependencies")
	q := u.Query()
	q.Set("dependencyId", dependencyID)
	u.RawQuery = q.Encode()

	c := dependency{}
	c.Description = description
	c.Consumer.DisplayName = consumerDisplayName
	if consumerOperationResourceName != "" && consumerExternalApiResourceName != "" {
		return nil, fmt.Errorf("consumerOperationResourceName and consumerExternalApiResourceName cannot be set together")
	}
	if consumerOperationResourceName != "" {
		c.Consumer.OperationResourceName = new(string)
		*c.Consumer.OperationResourceName = consumerOperationResourceName
	}
	if consumerExternalApiResourceName != "" {
		c.Consumer.ExternalApiResourceName = new(string)
		*c.Consumer.ExternalApiResourceName = consumerExternalApiResourceName
	}

	c.Supplier.DisplayName = supplierDisplayName
	if supplierOperationResourceName != "" && supplierExternalApiResourceName != "" {
		return nil, fmt.Errorf("supplierOperationResourceName and supplierExternalApiResourceName cannot be set together")
	}
	if supplierOperationResourceName != "" {
		c.Supplier.OperationResourceName = new(string)
		*c.Supplier.OperationResourceName = supplierOperationResourceName
	}
	if supplierExternalApiResourceName != "" {
		c.Supplier.ExternalApiResourceName = new(string)
		*c.Supplier.ExternalApiResourceName = supplierExternalApiResourceName
	}

	payload, err := json.Marshal(&c)
	if err != nil {
		return nil, err
	}
	respBody, err = apiclient.HttpClient(u.String(), string(payload))
	return respBody, err
}

func GetDependency(dependencyID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "dependencies", dependencyID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteDependency(dependencyID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "dependencies", dependencyID)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func ListDependencies(filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list("dependencies", filter, pageSize, pageToken)
}

func CreateDeployment(deploymentID string, displayName string, description string, deploymentName string,
	externalURI string, resourceURI string, endpoints []string, dep DeploymentType,
	env EnvironmentType, slo SloType,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "deployments")
	if deploymentID != "" {
		q := u.Query()
		q.Set("deploymentId", deploymentID)
		u.RawQuery = q.Encode()
	}

	payload, err := getDeployment(displayName, description, deploymentName, externalURI, resourceURI, endpoints, dep, env, slo)
	if err != nil {
		return nil, err
	}

	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func GetDeployment(deploymentID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "deployments", deploymentID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteDeployment(deploymentID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "deployments", deploymentID)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func ListDeployments(filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list("deployments", filter, pageSize, pageToken)
}

func UpdateDeployment(deploymentName string, displayName string, description string,
	externalURI string, resourceURI string, endpoints []string, dep DeploymentType,
	env EnvironmentType, slo SloType,
) (respBody []byte, err error) {
	updateMask := []string{}

	if displayName != "" {
		updateMask = append(updateMask, "displayName")
	}
	if description != "" {
		updateMask = append(updateMask, "description")
	}
	if externalURI != "" {
		updateMask = append(updateMask, "documentation")
	}
	if resourceURI != "" {
		updateMask = append(updateMask, "resource_uri")
	}
	if dep != "" {
		updateMask = append(updateMask, "deployment_type")
	}
	if env != "" {
		updateMask = append(updateMask, "environment")
	}
	if slo != "" {
		updateMask = append(updateMask, "slo")
	}
	if len(endpoints) > 0 {
		updateMask = append(updateMask, "endpoints")
	}
	if len(updateMask) == 0 {
		return nil, errors.New("Update mask is empty")
	}

	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "deployments", deploymentName)
	q := u.Query()
	q.Set("updateMask", strings.Join(updateMask, ","))
	u.RawQuery = q.Encode()

	payload, err := getDeployment(displayName, description, deploymentName,
		externalURI, resourceURI, endpoints, dep, env, slo)
	if err != nil {
		return nil, err
	}

	respBody, err = apiclient.HttpClient(u.String(), payload, "PATCH")
	return respBody, err
}

func getDeployment(displayName string, description string, deploymentName string,
	externalURI string, resourceURI string, endpoints []string, dep DeploymentType,
	env EnvironmentType, slo SloType,
) (string, error) {
	type documentation struct {
		ExternalURI string `json:"externalUri,omitempty"`
	}

	type attributeType struct {
		EnumValues enumValues `json:"enumValues,omitempty"`
	}

	type deployment struct {
		Name           string        `json:"name,omitempty"`
		DisplayName    string        `json:"displayName,omitempty"`
		Description    string        `json:"description,omitempty"`
		DeploymentType attributeType `json:"deploymentType,omitempty"`
		Documentation  documentation `json:"documentation,omitempty"`
		Environment    attributeType `json:"environment,omitempty"`
		Slo            attributeType `json:"slo,omitempty"`
		ResourceURI    string        `json:"resourceUri,omitempty"`
		Endpoints      []string      `json:"endpoints,omitempty"`
	}

	d := deployment{}

	if deploymentName != "" {
		d.Name = deploymentName
	}

	if displayName != "" {
		d.DisplayName = displayName
	}
	if description != "" {
		d.Description = description
	}
	if externalURI != "" {
		d.Documentation.ExternalURI = externalURI
	}
	if dep != "" {
		d.DeploymentType.EnumValues = getDeploymentEnum(dep).EnumValues
	}
	if env != "" {
		d.Environment.EnumValues = getEnvironmentEnum(env).EnumValues
	}
	if slo != "" {
		d.Slo.EnumValues = getSloEnum(slo).EnumValues
	}
	if resourceURI != "" {
		d.ResourceURI = resourceURI
	}
	if len(endpoints) > 0 {
		d.Endpoints = endpoints
	}

	payload, err := json.Marshal(&d)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func CreateExternalAPI(externalApiID string, displayName string, description string,
	endpoints []string, paths []string, externalUri string,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "externalApis")
	q := u.Query()
	q.Set("externalApiId", externalApiID)
	u.RawQuery = q.Encode()

	payload, err := getExternalApi(displayName, description, endpoints, paths, externalUri)
	if err != nil {
		return nil, err
	}
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func GetExternalAPI(externalApiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "externalApis", externalApiID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteExternalAPI(externalApiID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "externalApis", externalApiID)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func ListExternalAPIs(filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list("externalApis", filter, pageSize, pageToken)
}

func UpdateExternalAPI(externalApiID string, displayName string, description string,
	endpoints []string, paths []string, externalUri string,
) (respBody []byte, err error) {
	updateMask := []string{}

	if displayName != "" {
		updateMask = append(updateMask, "displayName")
	}
	if description != "" {
		updateMask = append(updateMask, "description")
	}
	if externalUri != "" {
		updateMask = append(updateMask, "documentation")
	}
	if len(paths) > 0 {
		updateMask = append(updateMask, "paths")
	}
	if len(endpoints) > 0 {
		updateMask = append(updateMask, "endpoints")
	}
	if len(updateMask) == 0 {
		return nil, errors.New("Update mask is empty")
	}

	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "externalApis", externalApiID)
	q := u.Query()
	q.Set("updateMask", strings.Join(updateMask, ","))
	u.RawQuery = q.Encode()

	payload, err := getExternalApi(displayName, description, endpoints, paths, externalUri)
	if err != nil {
		return nil, err
	}

	respBody, err = apiclient.HttpClient(u.String(), payload, "PATCH")
	return respBody, err
}

func getExternalApi(displayName string, description string,
	endpoints []string, paths []string, externalUri string,
) (string, error) {
	type documentation struct {
		ExternalURI string `json:"externalUri,omitempty"`
	}

	type extapi struct {
		DisplayName   string        `json:"displayName,omitempty"`
		Description   string        `json:"description,omitempty"`
		Documentation documentation `json:"documentation,omitempty"`
		Paths         []string      `json:"paths,omitempty"`
		Endpoints     []string      `json:"endpoints,omitempty"`
	}
	e := extapi{}
	if displayName != "" {
		e.DisplayName = displayName
	}
	if description != "" {
		e.Description = description
	}
	if externalUri != "" {
		e.Documentation.ExternalURI = externalUri
	}
	if len(paths) > 0 {
		e.Paths = paths
	}
	if len(endpoints) > 0 {
		e.Endpoints = endpoints
	}

	payload, err := json.Marshal(&e)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func CreateAttribute(attributeID string, displayName string, description string, scope string,
	dataType string, aValues []byte, cardinality int,
) (respBody []byte, err error) {
	type attributeScope string
	const (
		API           attributeScope = "API"
		VERSION       attributeScope = "VERSION"
		SPEC          attributeScope = "SPEC"
		API_OPERATION attributeScope = "API_OPERATION"
		DEPLOYMENT    attributeScope = "DEPLOYMENT"
		DEPENDENCY    attributeScope = "DEPENDENCY"
		DEFINITION    attributeScope = "DEFINITION"
		EXTERNAL_API  attributeScope = "EXTERNAL_API"
		PLUGIN        attributeScope = "PLUGIN"
	)

	type attributeDataType string
	const (
		ENUM   attributeDataType = "ENUM"
		JSON   attributeDataType = "JSON"
		STRING attributeDataType = "STRING"
	)

	type attribute struct {
		DisplayName   string            `json:"displayName,omitempty"`
		Description   string            `json:"description,omitempty"`
		Scope         attributeScope    `json:"scope,omitempty"`
		DataType      attributeDataType `json:"dataType,omitempty"`
		AllowedValues []allowedValue    `json:"allowedValues,omitempty"`
		Cardinality   int               `json:"cardinality,omitempty"`
	}

	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "attributes")
	q := u.Query()
	q.Set("attributeId", attributeID)
	u.RawQuery = q.Encode()

	a := attribute{}
	a.DisplayName = displayName
	a.Description = description
	a.Scope = attributeScope(scope)
	a.DataType = attributeDataType(dataType)
	a.Cardinality = cardinality

	if aValues != nil {
		var av []allowedValue
		err = json.Unmarshal(aValues, &av)
		if err != nil {
			return nil, err
		}
		a.AllowedValues = av
	}

	payload, err := json.Marshal(&a)
	if err != nil {
		return nil, err
	}
	respBody, err = apiclient.HttpClient(u.String(), string(payload))
	return respBody, err
}

func GetAttribute(attributeID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "attributes", attributeID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteAttribute(attributeID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "attributes", attributeID)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func ListAttributes(filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list("attributes", filter, pageSize, pageToken)
}

func list(resource string, filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, resource)
	q := u.Query()
	if filter != "" {
		q.Set("filter", filter)
	}
	if pageSize != -1 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func getSpecType(id string) enumAttributeValue {
	switch id {
	case "openapi":
		return getAttributeValues(allowedValue{
			Id:          "openapi",
			DisplayName: "OpenAPI Spec",
			Description: "OpenAPI Spec",
			Immutable:   true,
		})
	case "proto":
		return getAttributeValues(allowedValue{
			Id:          "proto",
			DisplayName: "Proto",
			Description: "Proto",
			Immutable:   true,
		})
	case "wsdl":
		return getAttributeValues(allowedValue{
			Id:          "wsdl",
			DisplayName: "WSDL",
			Description: "WSDL",
			Immutable:   true,
		})
	default:
		return getAttributeValues(allowedValue{
			Id:          "openapi",
			DisplayName: "OpenAPI Spec",
			Description: "OpenAPI Spec",
			Immutable:   true,
		})
	}
}

func getAttributeValues(a allowedValue) enumAttributeValue {
	l := []allowedValue{}
	l = append(l, a)
	e := enumValues{}
	e.Values = l

	s := enumAttributeValue{}
	s.EnumValues = enumValues{}
	s.EnumValues = e

	return s
}

func getDeploymentEnum(d DeploymentType) enumAttributeValue {
	switch d {
	case APIGEE:
		return getAttributeValues(allowedValue{
			Id:          "apigee",
			DisplayName: "Apigee",
			Description: "Apigee",
			Immutable:   true,
		})
	case APIGEE_HYBRID:
		return getAttributeValues(allowedValue{
			Id:          "apigee-hybrid",
			DisplayName: "Apigee Hybrid",
			Description: "Apigee Hybrid",
			Immutable:   true,
		})
	case APIGEE_EDGE_PRIVATE:
		return getAttributeValues(allowedValue{
			Id:          "apigee-edge-private",
			DisplayName: "Apigee Edge Private Cloud",
			Description: "Apigee Edge Private Cloud",
			Immutable:   true,
		})
	case APIGEE_EDGE_PUBLIC:
		return getAttributeValues(allowedValue{
			Id:          "apigee-edge-public",
			DisplayName: "Apigee Edge Public Cloud",
			Description: "Apigee Edge Public Cloud",
			Immutable:   true,
		})
	case MOCK_SERVER:
		return getAttributeValues(allowedValue{
			Id:          "mock-server",
			DisplayName: "Mock Server",
			Description: "Mock Server",
			Immutable:   true,
		})
	case CLOUD_API_GATEWAY:
		return getAttributeValues(allowedValue{
			Id:          "cloud-api-gateway",
			DisplayName: "Cloud API Gateway",
			Description: "Cloud API Gateway",
			Immutable:   true,
		})
	case CLOUD_ENDPOINTS:
		return getAttributeValues(allowedValue{
			Id:          "cloud-endpoints",
			DisplayName: "Cloud Endpoints",
			Description: "Cloud Endpoints",
			Immutable:   true,
		})
	case UNMANAGED:
		return getAttributeValues(allowedValue{
			Id:          "unmanaged",
			DisplayName: "Unmanaged",
			Description: "Unmanaged",
			Immutable:   true,
		})
	case OTHERS:
		return getAttributeValues(allowedValue{
			Id:          "others",
			DisplayName: "Others",
			Description: "Others",
			Immutable:   true,
		})
	default:
		return getAttributeValues(allowedValue{
			Id:          "others",
			DisplayName: "Others",
			Description: "Others",
			Immutable:   true,
		})
	}
}

func getEnvironmentEnum(e EnvironmentType) enumAttributeValue {
	switch e {
	case DEVELOPMENT:
		return getAttributeValues(allowedValue{
			Id:          "development",
			DisplayName: "Development",
			Description: "Development",
		})
	case STAGING:
		return getAttributeValues(allowedValue{
			Id:          "staging",
			DisplayName: "Staging",
			Description: "Staging",
		})
	case TEST:
		return getAttributeValues(allowedValue{
			Id:          "test",
			DisplayName: "Test",
			Description: "Test",
		})
	case PREPRODUCTION:
		return getAttributeValues(allowedValue{
			Id:          "pre-prod",
			DisplayName: "Pre-Production",
			Description: "Pre-Production",
		})
	case PRODUCTION:
		return getAttributeValues(allowedValue{
			Id:          "production",
			DisplayName: "Production",
			Description: "Production",
		})
	default:
		return getAttributeValues(allowedValue{
			Id:          "development",
			DisplayName: "Development",
			Description: "Development",
		})
	}
}

func getSloEnum(s SloType) enumAttributeValue {
	switch s {
	case SLO99_99:
		return getAttributeValues(allowedValue{
			Id:          "99-99",
			DisplayName: "99.99%",
			Description: "99.99% SLO",
		})
	case SLO99_95:
		return getAttributeValues(allowedValue{
			Id:          "99-95",
			DisplayName: "99.95%",
			Description: "99.95% SLO",
		})
	case SLO99_9:
		return getAttributeValues(allowedValue{
			Id:          "99-9",
			DisplayName: "99.9%",
			Description: "99.9% SLO",
		})
	case SLO99_5:
		return getAttributeValues(allowedValue{
			Id:          "99-5",
			DisplayName: "99.5%",
			Description: "99.5% SLO",
		})
	default:
		return getAttributeValues(allowedValue{
			Id:          "99-90",
			DisplayName: "99.90%",
			Description: "99.90% SLO",
		})
	}
}

func (d *DeploymentType) String() string {
	return string(*d)
}

func (d *DeploymentType) Set(r string) error {
	switch r {
	case "apigee", "apigee-hybrid", "apigee-edge-private", "apigee-edge-public", "mock-server", "cloud-api-gateway", "cloud-endpoints", "unmanaged", "others":
		*d = DeploymentType(r)
	default:
		return fmt.Errorf("must be one of %s, %s, %s, %s, %s, %s, %s, %s or %s",
			APIGEE, APIGEE_HYBRID, APIGEE_EDGE_PRIVATE, APIGEE_EDGE_PUBLIC, MOCK_SERVER, CLOUD_API_GATEWAY, CLOUD_ENDPOINTS, UNMANAGED, OTHERS)
	}
	return nil
}

func (d *DeploymentType) Type() string {
	return "deploymentType"
}

func (e *EnvironmentType) String() string {
	return string(*e)
}

func (e *EnvironmentType) Set(r string) error {
	switch r {
	case "apigee", "development", "test", "staging", "pre-prod", "prod":
		*e = EnvironmentType(r)
	default:
		return fmt.Errorf("must be one of %s, %s, %s, %s or %s",
			DEVELOPMENT, TEST, PRODUCTION, STAGING, PREPRODUCTION)
	}
	return nil
}

func (e *EnvironmentType) Type() string {
	return "environmentType"
}

func (s *SloType) String() string {
	return string(*s)
}

func (s *SloType) Set(r string) error {
	switch r {
	case "99-99", "99-95", "99-5", "99-9":
		*s = SloType(r)
	default:
		return fmt.Errorf("must be one of %s, %s, %s, or %s",
			SLO99_99, SLO99_95, SLO99_5, SLO99_9)
	}
	return nil
}

func (s *SloType) Type() string {
	return "sloType"
}

func getSpecIDList(s []byte) (sList []string) {

	type spec struct {
		Name         string                 `json:"name,omitempty"`
		DisplayName  string                 `json:"displayName,omitempty"`
		SpecType     map[string]interface{} `json:"specType,omitempty"`
		Details      map[string]interface{} `json:"details,omitempty"`
		LintResponse map[string]interface{} `json:"lintResponse,omitempty"`
	}

	type speclist struct {
		Specs []spec `json:"specs,omitempty"`
	}

	l := speclist{}

	if err := json.Unmarshal(s, &l); err != nil {
		clilog.Error.Println(err)
		return nil
	}

	for _, i := range l.Specs {
		sList = append(sList, filepath.Base(i.Name))
	}

	return sList
}

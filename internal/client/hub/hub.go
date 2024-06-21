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
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"internal/apiclient"
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
	DisplayName   string        `json:"displayName,omitempty"`
	SpecType      specType      `json:"specType,omitempty"`
	SourceURI     string        `json:"sourceUri,omitempty"`
	Contents      content       `json:"contents,omitempty"`
	Documentation documentation `json:"documentation,omitempty"`
	ParsingMode   string        `json:"parsingMode,omitempty"`
}

type documentation struct {
	ExternalUri string `json:"externalUri,omitempty"`
}

type content struct {
	Contents string `json:"contents,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

type allowedValues struct {
	Id          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
	Immutable   bool   `json:"immutable"`
}

type enumValues struct {
	Values []allowedValues `json:"values,omitempty"`
}

type specType struct {
	EnumValues enumValues `json:"enumValues,omitempty"`
}

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
	contents []byte, mimeType string, sourceURI string, documentation string) (respBody []byte, err error) {

	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs")
	q := u.Query()
	if specID != "" {
		q.Set("specId", specID)
	}
	u.RawQuery = q.Encode()

	s := spec{}
	s.DisplayName = displayName
	s.Documentation.ExternalUri = documentation
	s.Contents.Contents = base64.StdEncoding.EncodeToString(contents)

	if strings.Contains(mimeType, "yaml") || strings.Contains(mimeType, "yml") {
		s.Contents.MimeType = "application/yaml"
		s.SpecType = getAllowedValuesOpenAPI()
	} else if strings.Contains(mimeType, "json") {
		s.Contents.MimeType = "application/json"
		s.SpecType = getAllowedValuesOpenAPI()
	} else if strings.Contains(mimeType, "wsdl") {
		s.Contents.MimeType = "application/wsdl"
		s.SpecType = getAllowedValuesWSDL()
	} else if strings.Contains(mimeType, "proto") {
		s.Contents.MimeType = "application/text"
		s.SpecType = getAllowedValuesProto()
	} else {
		s.Contents.MimeType = "application/text"
	}
	s.SourceURI = sourceURI

	payload, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	respBody, err = apiclient.HttpClient(u.String(), string(payload))
	return respBody, err
}

func GetApiVersionsSpec(apiID string, versionID string, specID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs", specID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteApiVersionsSpec(apiID string, versionID string, specID string) (respBody []byte, err error) {
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

func ListApiVersionsSpec(apiID string, versionID string, filter string, pageSize int, pageToken string) (respBody []byte, err error) {
	return list(path.Join("apis", apiID, "versions", versionID, "specs"), filter, pageSize, pageToken)
}

func LintApiVersionSpec(apiID string, versionID string, specID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "apis", apiID, "versions", versionID, "specs", specID, ":lint")
	respBody, err = apiclient.HttpClient(u.String(), "")
	return respBody, err
}

func CreateDependency(dependencyID string, description string, consumerDisplayName string,
	consumerOperationResourceName string, consumerExternalApiResourceName string, supplierDisplayName string,
	supplierOperationResourceName string, supplierExternalApiResourceName string) (respBody []byte, err error) {

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

func CreateDeployment(deploymentID string, displayName string, description string,
	externalURI string, resourceURI string, endpoints []string, dep DeploymentType) (respBody []byte, err error) {

	type documentation struct {
		ExternalURI string `json:"externalUri,omitempty"`
	}

	type deploymentType struct {
		EnumValues enumValues `json:"enumValues,omitempty"`
	}

	type deployment struct {
		DisplayName    string         `json:"displayName,omitempty"`
		Description    string         `json:"description,omitempty"`
		DeploymentType deploymentType `json:"deploymentType,omitempty"`
		Documentation  documentation  `json:"documentation,omitempty"`
		ResourceURI    string         `json:"resourceUri,omitempty"`
		Endpoints      []string       `json:"endpoints,omitempty"`
	}

	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "deployments")
	q := u.Query()
	q.Set("deploymentId", deploymentID)
	u.RawQuery = q.Encode()

	d := deployment{}
	d.DisplayName = displayName
	d.Description = description
	d.Documentation.ExternalURI = externalURI
	d.ResourceURI = resourceURI
	d.Endpoints = endpoints

	payload, err := json.Marshal(&d)
	if err != nil {
		return nil, err
	}
	respBody, err = apiclient.HttpClient(u.String(), string(payload))
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

func CreateExternalAPI(externalApiId string, displayName string, description string,
	endpoints []string, paths []string, externalUri string) (respBody []byte, err error) {

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

	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "externalApi")
	q := u.Query()
	q.Set("externalApiId", externalApiId)
	u.RawQuery = q.Encode()

	e := extapi{}
	e.DisplayName = displayName
	e.Description = description
	e.Documentation.ExternalURI = externalUri
	e.Paths = paths
	e.Endpoints = endpoints

	payload, err := json.Marshal(&e)
	if err != nil {
		return nil, err
	}
	respBody, err = apiclient.HttpClient(u.String(), string(payload))
	return respBody, err
}

func GetExternalAPI(externalApiId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "externalApis", externalApiId)
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

func CreateAttribute(attributeID string, displayName string, description string, scope string,
	dataType string, aValues []byte, cardinality int) (respBody []byte, err error) {

	type attribute struct {
		DisplayName   string          `json:"displayName,omitempty"`
		Description   string          `json:"description,omitempty"`
		Scope         string          `json:"scope,omitempty"`
		DataType      string          `json:"dataType,omitempty"`
		AllowedValues []allowedValues `json:"allowedValues,omitempty"`
		Cardinality   int             `json:"cardinality,omitempty"`
	}
	u, _ := url.Parse(apiclient.GetApigeeRegistryURL())
	u.Path = path.Join(u.Path, "attributes")
	q := u.Query()
	q.Set("attributeId", attributeID)
	u.RawQuery = q.Encode()

	a := attribute{}
	a.DisplayName = displayName
	a.Description = description
	a.Scope = scope
	a.DataType = dataType
	a.Cardinality = cardinality

	if aValues != nil {
		var av []allowedValues
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

func getAllowedValuesOpenAPI() specType {
	a := allowedValues{
		Id:          "openapi",
		DisplayName: "OpenAPI Spec",
		Description: "OpenAPI Spec",
		Immutable:   true,
	}
	return getSpecType(a)
}

func getAllowedValuesProto() specType {
	a := allowedValues{
		Id:          "proto",
		DisplayName: "Proto",
		Description: "Proto",
		Immutable:   true,
	}
	return getSpecType(a)
}

func getAllowedValuesWSDL() specType {
	a := allowedValues{
		Id:          "wsdl",
		DisplayName: "WSDL",
		Description: "WSDL",
		Immutable:   true,
	}
	return getSpecType(a)
}

func getSpecType(a allowedValues) specType {
	l := []allowedValues{}
	l = append(l, a)
	e := enumValues{}
	e.Values = l

	s := specType{}
	s.EnumValues = enumValues{}
	s.EnumValues = e

	return s
}

func getDeploymentEnum(d DeploymentType) specType {
	switch d {
	case APIGEE:
		return getSpecType(allowedValues{
			Id:          "apigee",
			DisplayName: "Apigee",
			Description: "Apigee",
			Immutable:   true,
		})
	case APIGEE_HYBRID:
		return getSpecType(allowedValues{
			Id:          "apigee-hybrid",
			DisplayName: "Apigee Hybrid",
			Description: "Apigee Hybrid",
			Immutable:   true,
		})
	case APIGEE_EDGE_PRIVATE:
		return getSpecType(allowedValues{
			Id:          "apigee-edge-private",
			DisplayName: "Apigee Edge Private Cloud",
			Description: "Apigee Edge Private Cloud",
			Immutable:   true,
		})
	case APIGEE_EDGE_PUBLIC:
		return getSpecType(allowedValues{
			Id:          "apigee-edge-public",
			DisplayName: "Apigee Edge Public Cloud",
			Description: "Apigee Edge Public Cloud",
			Immutable:   true,
		})
	case MOCK_SERVER:
		return getSpecType(allowedValues{
			Id:          "mock-server",
			DisplayName: "Mock Server",
			Description: "Mock Server",
			Immutable:   true,
		})
	case CLOUD_API_GATEWAY:
		return getSpecType(allowedValues{
			Id:          "cloud-api-gateway",
			DisplayName: "Cloud API Gateway",
			Description: "Cloud API Gateway",
			Immutable:   true,
		})
	case CLOUD_ENDPOINTS:
		return getSpecType(allowedValues{
			Id:          "cloud-endpoints",
			DisplayName: "Cloud Endpoints",
			Description: "Cloud Endpoints",
			Immutable:   true,
		})
	case UNMANAGED:
		return getSpecType(allowedValues{
			Id:          "unmanaged",
			DisplayName: "Unmanaged",
			Description: "Unmanaged",
			Immutable:   true,
		})
	case OTHERS:
		return getSpecType(allowedValues{
			Id:          "others",
			DisplayName: "Others",
			Description: "Others",
			Immutable:   true,
		})
	default:
		return getSpecType(allowedValues{
			Id:          "others",
			DisplayName: "Others",
			Description: "Others",
			Immutable:   true,
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

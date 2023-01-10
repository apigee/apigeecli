// Copyright 2022 Google LLC
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

package bundlegen

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/apigee/apigeecli/apiclient"
	apiproxy "github.com/apigee/apigeecli/bundlegen/apiproxydef"
	"github.com/apigee/apigeecli/bundlegen/policies"
	"github.com/apigee/apigeecli/bundlegen/proxies"
	"github.com/apigee/apigeecli/bundlegen/targets"
	"github.com/apigee/apigeecli/clilog"
	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/ghodss/yaml"
)

type backendDef struct {
	Address         string
	JwtAudience     string
	DisableAuth     bool
	PathTranslation string
	Deadline        int //seconds
}

type googleManagementDef struct {
	Metrics interface{}
	Quota   googleLimitsDef
}

type googleLimitsDef struct {
	Limits []googleLimitDef
}

type googleLimitDef struct {
	Name   string
	Metric string
	Unit   string
	Value  ValueDef
}

type ValueDef struct {
	Standard string
}

var doc2 *openapi2.T

const NoAuthTargetName = "default"
const GoogleAuthTargetName = "google-auth"

var amPolicyContent = make(map[string]string)

var allowValue, apiName string

var googMgmt googleManagementDef
var quotaList []quotaDef
var defaultBackend backendDef

func LoadSwaggerFromUri(endpoint string) (string, []byte, error) {
	var docType string

	u, err := url.Parse(endpoint)
	clilog.Info.Printf("%v\n", u)
	if err != nil {
		clilog.Error.Println(err)
		return "", nil, err
	}

	name := path.Base(u.Path)
	if docType = filepath.Ext(name); docType == "" {
		docType = "json"
		name = name + ".json"
	}
	clilog.Info.Printf("docType: %s\n", docType)

	if err = apiclient.DownloadResource(endpoint, name, docType, false); err != nil {
		clilog.Error.Println(err)
		return "", nil, err
	}
	defer os.Remove(name)
	return LoadSwaggerFromFile(name)
}

func LoadSwaggerFromFile(filePath string) (string, []byte, error) {
	var err error
	var jsonContent, swaggerBytes, swaggerJsonBytes []byte

	if swaggerBytes, err = utils.ReadFile(filePath); err != nil {
		clilog.Error.Println(err)
		return "", nil, err
	}

	//convert yaml to json
	if isFileYaml(filePath) {
		if swaggerJsonBytes, err = yaml.YAMLToJSON(swaggerBytes); err != nil {
			clilog.Error.Println(err)
			return "", nil, err
		}
		swaggerBytes = swaggerJsonBytes
	}

	if err = json.Unmarshal(swaggerBytes, &doc2); err != nil {
		clilog.Error.Println(err)
		return "", nil, err
	}

	if jsonContent, err = doc2.MarshalJSON(); err != nil {
		clilog.Error.Println(err)
		return "", nil, err
	}

	clilog.Info.Printf("%s", string(jsonContent))

	return filepath.Base(filePath), jsonContent, err
}

func GenerateAPIProxyFromSwagger(name string,
	oasDocName string,
	basePath string,
	addCORS bool) (string, error) {

	var err error

	//load the security definitions
	loadSwaggerSecurityRequirements(doc2.SecurityDefinitions)

	//load google extensions
	err = loadGoogleExtensions()
	if err != nil {
		clilog.Error.Println(err)
		return name, err
	}

	if name != "" {
		apiproxy.SetDisplayName(name)
		//set the name for use when generating the bundle
		apiName = name
	} else if apiName != "" {
		apiproxy.SetDisplayName(apiName)
	} else {
		return name, fmt.Errorf("neither x-google-api-name nor name was set")
	}

	if doc2.Info.Description != "" {
		apiproxy.SetDescription(doc2.Info.Description)
	}

	apiproxy.SetCreatedAt()
	apiproxy.SetLastModifiedAt()
	apiproxy.SetConfigurationVersion()
	apiproxy.AddProxyEndpoint("default")

	if doc2.BasePath == "" {
		return name, fmt.Errorf("basePath is missing from the Swagger file. Please add a basePath")
	}

	apiproxy.SetBasePath(doc2.BasePath)
	proxies.NewProxyEndpoint(doc2.BasePath, true)

	//add global security policies
	if securityScheme := getSwaggerSecurityRequirements(doc2.Security); securityScheme.SchemeName != "" {
		if securityScheme.APIKeyPolicy.APIKeyPolicyEnabled {
			proxies.AddStepToPreFlowRequest("Verify-API-Key-" + securityScheme.SchemeName)
			enableSecurityPolicy(securityScheme.SchemeName, "apikey")
		} else if securityScheme.JWTPolicy.JWTPolicyEnabled {
			proxies.AddStepToPreFlowRequest("VerifyJWT-" + securityScheme.SchemeName)
			enableSecurityPolicy(securityScheme.SchemeName, "jwt")
		}
	}

	if err = generateSwaggerFlows(doc2.Paths); err != nil {
		clilog.Error.Println(err)
		return name, err
	}

	//handle unhandled requests
	if allowValue == "configured" {
		proxies.AddFlow("Unknown Request", "", "", "Handle unknown requests")
		proxies.AddStepToFlowRequest("Raise-Fault-Unknown-Request", "Unknown Request")
		apiproxy.AddPolicy("Raise-Fault-Unknown-Request")
	}

	if defaultBackend.Address != "" { //there is a default address
		if err = addBackend(defaultBackend); err != nil {
			return name, err
		}
		if defaultBackend.JwtAudience != "" {
			proxies.AddStepToPreFlowRequest("Copy-Auth-Var")
			apiproxy.AddPolicy("Copy-Auth-Var")
			policies.EnableCopyAuthPolicy()
		}
	}

	for _, securityScheme := range securitySchemesList.SecuritySchemes {
		if securityScheme.JWTPolicy.JWTPolicyEnabled {
			apiproxy.AddPolicy("VerifyJWT-" + securityScheme.SchemeName)
		} else if securityScheme.APIKeyPolicy.APIKeyPolicyEnabled {
			apiproxy.AddPolicy("Verify-API-Key-" + securityScheme.SchemeName)
		}
	}

	if addCORS {
		proxies.AddStepToPreFlowRequest("Add-CORS")
		apiproxy.AddPolicy("Add-CORS")
	}

	return name, nil
}

func loadSecurityDefinition(secDefName string, securityScheme openapi2.SecurityScheme) (secScheme securitySchemesDef) {
	secScheme = securitySchemesDef{}
	jwtPolicy := jwtPolicyDef{}
	apiKeyPolicy := apiKeyPolicyDef{}

	if securityScheme.Type == "oauth2" {
		secScheme.SchemeName = secDefName

		for extensionName, extensionValue := range securityScheme.Extensions {
			if extensionName == "x-google-audiences" {
				jwtPolicy.Audience = strings.ReplaceAll(fmt.Sprintf("%s", extensionValue), "\"", "")
			}
			if extensionName == "x-google-issuer" {
				jwtPolicy.Issuer = strings.ReplaceAll(fmt.Sprintf("%s", extensionValue), "\"", "")
			}
			if extensionName == "x-google-jwks_uri" {
				jwtPolicy.JwkUri = strings.ReplaceAll(fmt.Sprintf("%s", extensionValue), "\"", "")
			}
			if extensionName == "x-google-jwt-locations" {
				locations := []map[string]string{}
				str := fmt.Sprintf("%s", extensionValue)

				if err := json.Unmarshal([]byte(str), &locations); err != nil {
					clilog.Error.Println(err)
				}
				if len(locations) > 0 {
					//deal with only the first location
					jwtPolicy.Location = make(map[string]string)
					for locationKey, locationValue := range locations[0] {
						jwtPolicy.Location[locationKey] = locationValue
						jwtPolicy.Source = "tokenVar"
					}
				}
			}
		}
		secScheme.JWTPolicy = jwtPolicy
	} else if securityScheme.Type == "apiKey" {
		secScheme.SchemeName = secDefName
		apiKeyPolicy.APIKeyPolicyEnabled = true
		apiKeyPolicy.APIKeyLocation = securityScheme.In
		apiKeyPolicy.APIKeyName = securityScheme.Name
		secScheme.APIKeyPolicy = apiKeyPolicy
	}

	return secScheme
}

func loadGoogleExtensions() (err error) {

	for extensionName, extensionValue := range doc2.Extensions {
		clilog.Info.Printf("Found extension: %s", extensionName)
		if extensionName == "x-google-management" {
			if err := parseManagementExtension(extensionValue); err != nil {
				return err
			}
		} else if extensionName == "x-google-allow" {
			allowValue := strings.ReplaceAll(fmt.Sprintf("%s", extensionValue), "\"", "")
			clilog.Info.Printf("Allow Value: %s\n", allowValue)
			if allowValue != "configured" && allowValue != "all" {
				return fmt.Errorf("invalid value for x-google-allow: %s", allowValue)
			}
		} else if extensionName == "x-google-api-name" {
			clilog.Info.Printf("Found API Name: %s\n", extensionValue)
			if apiName, err = parseApiExtension(extensionValue); err != nil {
				return err
			}
		} else if extensionName == "x-google-backend" {
			if defaultBackend, err = parseBackendExtension(extensionValue, false); err != nil {
				return err
			}
			clilog.Info.Printf("Found default backend: %v", defaultBackend)
		}
	}
	return nil
}

func enableSecurityPolicy(name string, policyType string) {
	for index, securityScheme := range securitySchemesList.SecuritySchemes {
		if securityScheme.SchemeName == name {
			if policyType == "jwt" {
				securitySchemesList.SecuritySchemes[index].JWTPolicy.JWTPolicyEnabled = true
			} else if policyType == "apikey" {
				securitySchemesList.SecuritySchemes[index].APIKeyPolicy.APIKeyPolicyEnabled = true
			}
		}
	}
}

func getSecurityType(secName string) securitySchemesDef {
	for _, securityScheme := range securitySchemesList.SecuritySchemes {
		if securityScheme.SchemeName == secName {
			return securityScheme
		}
	}
	return securitySchemesDef{}
}

func getSwaggerSecurityRequirements(securityRequirements openapi2.SecurityRequirements) securitySchemesDef {
	for _, secReq := range securityRequirements {
		for secReqName := range secReq {
			return getSecurityType(secReqName)
		}
	}
	return securitySchemesDef{}
}

func loadSwaggerSecurityRequirements(securityDefinitions map[string]*openapi2.SecurityScheme) {
	for secDefName, secDef := range securityDefinitions {
		clilog.Info.Printf("Loading Security Definition: %s\n", secDefName)
		securitySchemesList.SecuritySchemes = append(securitySchemesList.SecuritySchemes, loadSecurityDefinition(secDefName, *secDef))
	}
}

func getSwaggerHTTPMethod(pathItem openapi2.PathItem, keyPath string) (map[string]pathDetailDef, error) {

	var err error
	pathMap := make(map[string]pathDetailDef)
	alternateOperationId := strings.ReplaceAll(keyPath, "\\", "_")

	if pathItem.Get != nil {
		getPathDetail := pathDetailDef{}
		getPathDetail.Verb = "GET"
		getPathDetail.Path = keyPath
		if pathItem.Get.OperationID != "" {
			getPathDetail.OperationID = pathItem.Get.OperationID
		} else {
			getPathDetail.OperationID = "get_" + alternateOperationId
		}
		if pathItem.Get.Description != "" {
			getPathDetail.Description = pathItem.Get.Description
		}
		if pathItem.Get.Security != nil {
			securityRequirements := openapi2.SecurityRequirements(*pathItem.Get.Security)
			getPathDetail.SecurityScheme = getSwaggerSecurityRequirements(securityRequirements)
			if getPathDetail.SecurityScheme.JWTPolicy.Audience != "" {
				getPathDetail.SecurityScheme.JWTPolicy.JWTPolicyEnabled = true
			}
		}
		//check for google extensions
		if pathItem.Get.Extensions != nil {
			if getPathDetail, err = processPathSwaggerExtensions(pathItem.Get.Extensions, getPathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["get"] = getPathDetail
	}

	if pathItem.Post != nil {
		postPathDetail := pathDetailDef{}
		postPathDetail.Verb = "POST"
		postPathDetail.Path = keyPath
		if pathItem.Post.OperationID != "" {
			postPathDetail.OperationID = pathItem.Post.OperationID
		} else {
			postPathDetail.OperationID = "post_" + alternateOperationId
		}
		if pathItem.Post.Description != "" {
			postPathDetail.Description = pathItem.Post.Description
		}
		if pathItem.Post.Security != nil {
			securityRequirements := openapi2.SecurityRequirements(*pathItem.Post.Security)
			postPathDetail.SecurityScheme = getSwaggerSecurityRequirements(securityRequirements)
			if postPathDetail.SecurityScheme.JWTPolicy.Audience != "" {
				postPathDetail.SecurityScheme.JWTPolicy.JWTPolicyEnabled = true
			}
		}
		//check for google extensions
		if pathItem.Post.Extensions != nil {
			if postPathDetail, err = processPathSwaggerExtensions(pathItem.Post.Extensions, postPathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["post"] = postPathDetail
	}

	if pathItem.Put != nil {
		putPathDetail := pathDetailDef{}
		putPathDetail.Verb = "PUT"
		putPathDetail.Path = keyPath
		if pathItem.Put.OperationID != "" {
			putPathDetail.OperationID = pathItem.Put.OperationID
		} else {
			putPathDetail.OperationID = "put_" + alternateOperationId
		}
		if pathItem.Put.Description != "" {
			putPathDetail.Description = pathItem.Put.Description
		}
		if pathItem.Put.Security != nil {
			securityRequirements := openapi2.SecurityRequirements(*pathItem.Put.Security)
			putPathDetail.SecurityScheme = getSwaggerSecurityRequirements(securityRequirements)
			if putPathDetail.SecurityScheme.JWTPolicy.Audience != "" {
				putPathDetail.SecurityScheme.JWTPolicy.JWTPolicyEnabled = true
			}
		}
		//check for google extensions
		if pathItem.Put.Extensions != nil {
			if putPathDetail, err = processPathSwaggerExtensions(pathItem.Put.Extensions, putPathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["put"] = putPathDetail
	}

	if pathItem.Patch != nil {
		patchPathDetail := pathDetailDef{}
		patchPathDetail.Verb = "PATCH"
		patchPathDetail.Path = keyPath
		if pathItem.Patch.OperationID != "" {
			patchPathDetail.OperationID = pathItem.Patch.OperationID
		} else {
			patchPathDetail.OperationID = "patch_" + alternateOperationId
		}
		if pathItem.Patch.Description != "" {
			patchPathDetail.Description = pathItem.Patch.Description
		}
		if pathItem.Patch.Security != nil {
			securityRequirements := openapi2.SecurityRequirements(*pathItem.Patch.Security)
			patchPathDetail.SecurityScheme = getSwaggerSecurityRequirements(securityRequirements)
			if patchPathDetail.SecurityScheme.JWTPolicy.Audience != "" {
				patchPathDetail.SecurityScheme.JWTPolicy.JWTPolicyEnabled = true
			}
		}
		//check for google extensions
		if pathItem.Patch.Extensions != nil {
			if patchPathDetail, err = processPathSwaggerExtensions(pathItem.Patch.Extensions, patchPathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["patch"] = patchPathDetail
	}

	if pathItem.Delete != nil {
		deletePathDetail := pathDetailDef{}
		deletePathDetail.Verb = "DELETE"
		deletePathDetail.Path = keyPath
		if pathItem.Delete.OperationID != "" {
			deletePathDetail.OperationID = pathItem.Delete.OperationID
		} else {
			deletePathDetail.OperationID = "delete_" + alternateOperationId
		}
		if pathItem.Delete.Description != "" {
			deletePathDetail.Description = pathItem.Delete.Description
		}
		if pathItem.Delete.Security != nil {
			securityRequirements := openapi2.SecurityRequirements(*pathItem.Delete.Security)
			deletePathDetail.SecurityScheme = getSwaggerSecurityRequirements(securityRequirements)
			if deletePathDetail.SecurityScheme.JWTPolicy.Audience != "" {
				deletePathDetail.SecurityScheme.JWTPolicy.JWTPolicyEnabled = true
			}
		}
		//check for google extensions
		if pathItem.Delete.Extensions != nil {
			if deletePathDetail, err = processPathSwaggerExtensions(pathItem.Delete.Extensions, deletePathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["delete"] = deletePathDetail
	}

	if pathItem.Options != nil {
		optionsPathDetail := pathDetailDef{}
		if pathItem.Options.OperationID != "" {
			optionsPathDetail.OperationID = pathItem.Options.OperationID
		} else {
			optionsPathDetail.OperationID = "options_" + alternateOperationId
		}
		if pathItem.Options.Description != "" {
			optionsPathDetail.Description = pathItem.Options.Description
		}
		pathMap["options"] = optionsPathDetail
	}

	if pathItem.Head != nil {
		headPathDetail := pathDetailDef{}
		if pathItem.Head.OperationID != "" {
			headPathDetail.OperationID = pathItem.Head.OperationID
		} else {
			headPathDetail.OperationID = "head_" + alternateOperationId
		}
		if pathItem.Head.Description != "" {
			headPathDetail.Description = pathItem.Head.Description
		}
		pathMap["head"] = headPathDetail
	}

	return pathMap, nil
}

func generateSwaggerFlows(paths map[string]*openapi2.PathItem) (err error) {
	for pathName, pathItem := range paths {
		pathMap, err := getSwaggerHTTPMethod(*pathItem, pathName)
		if err != nil {
			return err
		}
		for method, pathDetail := range pathMap {
			if !proxies.FlowExists(pathDetail.OperationID) {
				proxies.AddFlow(pathDetail.OperationID, replacePathWithWildCard(pathName), method, pathDetail.Description)

				if pathDetail.Backend != (backendDef{}) {
					if pathDetail.Backend.JwtAudience != "" {
						if !targets.IsExists(GoogleAuthTargetName) {
							targets.NewTargetEndpoint(GoogleAuthTargetName, pathDetail.Backend.Address, "", pathDetail.Backend.JwtAudience, "")
						}
						if err = targets.AddFlow(GoogleAuthTargetName, pathDetail.OperationID, replacePathWithWildCard(pathName), method, pathDetail.Description); err != nil {
							return err
						}
					} else {
						if !targets.IsExists(NoAuthTargetName) {
							targets.NewTargetEndpoint(NoAuthTargetName, pathDetail.Backend.Address, "", "", "")
						}
						if err = targets.AddFlow(NoAuthTargetName, pathDetail.OperationID, replacePathWithWildCard(pathName), method, pathDetail.Description); err != nil {
							return err
						}
					}
				}
			}
			if pathDetail.AssignMessage != "" {
				if pathDetail.Backend != (backendDef{}) {
					if pathDetail.Backend.JwtAudience != "" {
						if err = targets.AddStepToFlowRequest(GoogleAuthTargetName, "AM-"+pathDetail.OperationID, pathDetail.OperationID); err != nil {
							return err
						}
					} else {
						if err = targets.AddStepToFlowRequest(NoAuthTargetName, "AM-"+pathDetail.OperationID, pathDetail.OperationID); err != nil {
							return err
						}
					}
					apiproxy.AddPolicy("AM-" + pathDetail.OperationID)
				}
			}
			if pathDetail.SecurityScheme.JWTPolicy.JWTPolicyEnabled {
				//handle jwt locations
				if len(pathDetail.SecurityScheme.JWTPolicy.Location) != 0 { //jwt-location is specified
					if err = proxies.AddStepToFlowRequest("ExtractJWT-"+pathDetail.SecurityScheme.SchemeName, pathDetail.OperationID); err != nil {
						return err
					}
					apiproxy.AddPolicy("ExtractJWT-" + pathDetail.SecurityScheme.SchemeName)
				}
				//end handle jwt locations
				if err = proxies.AddStepToFlowRequest("VerifyJWT-"+pathDetail.SecurityScheme.SchemeName, pathDetail.OperationID); err != nil {
					return err
				}
				apiproxy.AddPolicy("VerifyJWT-" + pathDetail.SecurityScheme.SchemeName)
				enableSecurityPolicy(pathDetail.SecurityScheme.SchemeName, "jwt")
				//copy the original authorization header to X-Forwarded-Authorization
				// source: https://cloud.google.com/endpoints/docs/openapi/openapi-extensions#jwt_audience
				if err = proxies.AddStepToFlowRequest("Copy-Auth-Var", pathDetail.OperationID); err != nil {
					return err
				}
				apiproxy.AddPolicy("Copy-Auth-Var")
				policies.EnableCopyAuthPolicy()
			}
			if pathDetail.SecurityScheme.APIKeyPolicy.APIKeyPolicyEnabled {
				if err = proxies.AddStepToFlowRequest("Verify-API-Key-"+pathDetail.SecurityScheme.SchemeName, pathDetail.OperationID); err != nil {
					return err
				}
			}
			if pathDetail.Quota.QuotaEnabled {
				if err = proxies.AddStepToFlowRequest("Quota-"+pathDetail.Quota.QuotaName, pathDetail.OperationID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func parseBackendExtension(i interface{}, operation bool) (backendDef, error) {

	backend := backendDef{}
	var jsonMap map[string]interface{}
	var disableAuth string

	str := fmt.Sprintf("%s", i)

	if err := json.Unmarshal([]byte(str), &jsonMap); err != nil {
		return backendDef{}, err
	}

	if jsonMap["address"] != "" {
		tmp := fmt.Sprintf("%s", jsonMap["address"])
		tmp = strings.Replace(tmp, "<nil>", "", -1)
		if tmp != "" {
			backend.Address = tmp
		}
	}

	if jsonMap["disable_auth"] != "" {
		disableAuth = fmt.Sprintf("%v", jsonMap["disable_auth"])
		disableAuth = strings.Replace(disableAuth, "<nil>", "", -1)
		if disableAuth != "" {
			backend.DisableAuth, _ = strconv.ParseBool(disableAuth)
		}
	}

	//If address is not set, ESPv2 will automatically set disable_auth to true
	if backend.Address == "" {
		clilog.Info.Println("Address not set, disabling auth")
		backend.DisableAuth = true
	}

	if jsonMap["jwt_audience"] != "" {
		tmp := fmt.Sprintf("%#v", jsonMap["jwt_audience"])
		tmp = strings.Replace(tmp, "<nil>", "", -1)
		if tmp != "" {
			backend.JwtAudience = tmp
		}
	}

	if backend.JwtAudience != "" && disableAuth != "" {
		return backend, fmt.Errorf("both jwt_audience and disable_auth cannot be set")
	}

	//If an operation uses x-google-backend but does not specify either jwt_audience
	// or disable_auth, ESPv2 will automatically default the jwt_audience to match the address
	clilog.Info.Printf("Operation: %t, Audience %s, disable_auth: %s\n", operation, backend.JwtAudience, disableAuth)
	if operation && backend.JwtAudience == "" && disableAuth == "" {
		backend.JwtAudience = backend.Address
	}

	if jsonMap["deadline"] != "" {
		tmp := fmt.Sprintf("%v", jsonMap["deadline"])
		tmp = strings.Replace(tmp, "<nil>", "", -1)
		if tmp != "" {
			backend.Deadline, _ = strconv.Atoi(tmp)
		}
	}

	if jsonMap["path_transalation"] != "" {
		tmp := fmt.Sprintf("%v", jsonMap["path_transalation"])
		tmp = strings.Replace(tmp, "<nil>", "", -1)
		if tmp != "" && tmp != "APPEND_PATH_TO_ADDRESS" && tmp != "CONSTANT_ADDRESS" {
			return backend, fmt.Errorf("invalid path translation options: %s", tmp)
		} else {
			backend.PathTranslation = tmp
		}
	} else {
		backend.PathTranslation = "APPEND_PATH_TO_ADDRESS"
	}

	clilog.Info.Printf("Parsed Backend: %v\n", backend)

	return backend, nil
}

func parseApiExtension(i interface{}) (string, error) {
	str := strings.ReplaceAll(fmt.Sprintf("%s", i), "\"", "")
	if str == "" {
		return "", fmt.Errorf("x-google-api-name not found")
	}
	return str, nil
}

func parseManagementExtension(i interface{}) error {
	googMgmt = googleManagementDef{}
	str := fmt.Sprintf("%s", i)

	clilog.Info.Printf("Raw x-google-management: %s\n", str)

	if err := json.Unmarshal([]byte(str), &googMgmt); err != nil {
		return err
	}

	for _, limit := range googMgmt.Quota.Limits {
		quota := quotaDef{}
		quota.QuotaName = limit.Metric
		quota.QuotaTimeUnitLiteral = "minute"
		quota.QuotaIntervalLiteral = "1"
		quota.QuotaAllowLiteral = limit.Value.Standard
		clilog.Info.Printf("Found quota definition: %v\n", quota)
		quotaList = append(quotaList, quota)
	}
	return nil
}

func parseQuotaExtension(i interface{}) (quotaDef, error) {
	var jsonMap map[string]interface{}

	str := fmt.Sprintf("%s", i)

	if err := json.Unmarshal([]byte(str), &jsonMap); err != nil {
		return quotaDef{}, err
	}

	//search defined quota
	for name := range jsonMap {
		tmp := fmt.Sprintf("%v", jsonMap[name])
		tmp = strings.ReplaceAll(tmp, "map[", "")
		tmp = strings.ReplaceAll(tmp, "]", "")
		keyValue := strings.Split(tmp, ":")
		for index, quota := range quotaList {
			if keyValue[0] == quota.QuotaName {
				quotaList[index].QuotaAllowLiteral = keyValue[1]
				quotaList[index].QuotaEnabled = true
				quotaList[index].QuotaIdentiferLiteral = "organization.name" //this mimics rate limit per project which endpoints does.
				//store the XML policy contents
				quotaPolicyContent[quota.QuotaName] = policies.AddQuotaPolicy("Quota-"+quotaList[index].QuotaName,
					quotaList[index].QuotaConfigStepName,
					quotaList[index].QuotaAllowRef,
					quotaList[index].QuotaAllowLiteral,
					quotaList[index].QuotaIntervalRef,
					quotaList[index].QuotaIntervalLiteral,
					quotaList[index].QuotaTimeUnitRef,
					quotaList[index].QuotaTimeUnitLiteral,
					quotaList[index].QuotaIdentifierRef,
					quotaList[index].QuotaIdentiferLiteral)
				return quotaList[index], nil
			}
		}
	}
	return quotaDef{}, fmt.Errorf("quota defined in the path did not match the definition in x-google-management")
}

func getConditionString(matchespath string, verb string) string {
	re := regexp.MustCompile(`{\w+}`)
	if ok := re.Match([]byte(matchespath)); ok {
		matchespath = re.ReplaceAllString(matchespath, "*")
	}
	return "(proxy.pathsuffix MatchesPath \"" + matchespath + "\") and (request.verb = \"" + strings.ToUpper(verb) + "\")"
}

func processPathSwaggerExtensions(extensions map[string]interface{}, pathDetail pathDetailDef) (pathDetailDef, error) {
	var err error
	for extensionName, extensionValue := range extensions {
		if extensionName == "x-google-backend" {
			//process google-backed
			backend, err := parseBackendExtension(extensionValue, true)
			if err != nil {
				return pathDetail, err
			}
			if backend.JwtAudience != "" {
				proxies.AddRoute(pathDetail.OperationID, GoogleAuthTargetName, getConditionString(pathDetail.Path, pathDetail.Verb))
			} else {
				proxies.AddRoute(pathDetail.OperationID, NoAuthTargetName, getConditionString(pathDetail.Path, pathDetail.Verb))
			}
			pathDetail.AssignMessage = policies.AddSetTargetEndpoint("AM-"+pathDetail.OperationID, backend.Address, backend.PathTranslation)
			setAMPolicy(pathDetail.OperationID, pathDetail.AssignMessage)
			if err = addBackend(backend); err != nil {
				return pathDetail, err
			}
			pathDetail.Backend = backend
		} else if extensionName == "x-google-quota" {
			//process quota
			quota, err := parseQuotaExtension(extensionValue)
			if err != nil {
				return pathDetail, err
			}
			pathDetail.Quota = quota
		}
	}
	return pathDetail, err
}

func GetAMPolicies() map[string]string {
	return amPolicyContent
}

func GetGoogleApiName() string {
	return apiName
}

func GetAllowDefinition() string {
	return allowValue
}

func setAMPolicy(name string, content string) {
	amPolicyContent[name] = content
}

func addBackend(backend backendDef) (err error) {
	if backend.Address == "" {
		return fmt.Errorf("address is a mandatory field in x-google-backend")
	}
	//if there is a jwt_audience specified and auth is not disabled, use google auth
	clilog.Info.Printf("JwtAudience %s and DisableAuth %t\n", backend.JwtAudience, backend.DisableAuth)
	if backend.JwtAudience != "" && !backend.DisableAuth {
		if !targets.IsExists(GoogleAuthTargetName) {
			clilog.Info.Println("Adding Google Auth Target Server")
			apiproxy.AddTargetEndpoint(GoogleAuthTargetName)
			targets.NewTargetEndpoint(GoogleAuthTargetName, backend.Address, "", backend.JwtAudience, "")
			//at the moment one cannot have different deadlines per target.
			if backend.Deadline > 0 {
				targets.AddTargetEndpointProperty(GoogleAuthTargetName, "connect.timeout.millis", fmt.Sprintf("%d", backend.Deadline*1000))
			}
		} else {
			clilog.Info.Println("Google Auth Target Server already exists")
		}
	} else {
		if !targets.IsExists(NoAuthTargetName) {
			clilog.Info.Println("Adding Default Target Server")
			apiproxy.AddTargetEndpoint(NoAuthTargetName)
			targets.NewTargetEndpoint(NoAuthTargetName, backend.Address, "", "", "")
			//at the moment one cannot have different deadlines per target.
			if backend.Deadline > 0 {
				targets.AddTargetEndpointProperty(NoAuthTargetName, "connect.timeout.millis", fmt.Sprintf("%d", backend.Deadline*1000))
			}
		} else {
			clilog.Info.Println("Default Target Server already exists")
		}
	}
	return nil
}

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

package bundlegen

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	apiproxy "github.com/srinandan/apigeecli/bundlegen/apiproxydef"
	proxies "github.com/srinandan/apigeecli/bundlegen/proxies"
	targets "github.com/srinandan/apigeecli/bundlegen/targets"
)

type pathDetailDef struct {
	OperationID    string
	Description    string
	SecurityScheme securitySchemesDef
}

type oAuthPolicyDef struct {
	OAuthPolicyEnabled bool
}

type securitySchemesListDef struct {
	SecuritySchemes []securitySchemesDef
}

type securitySchemesDef struct {
	SchemeName   string
	OAuthPolicy  oAuthPolicyDef
	APIKeyPolicy apiKeyPolicyDef
}

type apiKeyPolicyDef struct {
	APIKeyPolicyEnabled bool
	APIKeyLocation      string
	APIKeyName          string
}

var generateSetTarget bool

var securitySchemesList = securitySchemesListDef{}

var doc *openapi3.T

func LoadDocumentFromFile(filePath string, validate bool, formatValidation bool) (string, []byte, error) {
	var err error
	var jsonContent []byte

	doc, err = openapi3.NewLoader().LoadFromFile(filePath)
	if err != nil {
		return "", nil, err
	}

	//add custom string definitions
	openapi3.DefineStringFormat("uuid", openapi3.FormatOfStringForUUIDOfRFC4122)

	if !formatValidation {
		openapi3.SchemaFormatValidationDisabled = true
	}

	if validate {
		if err = doc.Validate(openapi3.NewLoader().Context); err != nil {
			return "", nil, err
		}
	}

	if jsonContent, err = doc.MarshalJSON(); err != nil {
		return "", nil, err
	}

	if isFileYaml(filePath) {
		yamlContent, err := yaml.JSONToYAML(jsonContent)
		return filepath.Base(filePath), yamlContent, err
	} else {
		return filepath.Base(filePath), jsonContent, err
	}
}

func LoadDocumentFromURI(uri string, validate bool, formatValidation bool) (string, []byte, error) {
	var err error
	var jsonContent []byte

	u, err := url.Parse(uri)
	if err != nil {
		return "", nil, err
	}

	doc, err = openapi3.NewLoader().LoadFromURI(u)
	if err != nil {
		return "", nil, err
	}

	//add custom string definitions
	openapi3.DefineStringFormat("uuid", openapi3.FormatOfStringForUUIDOfRFC4122)

	if !formatValidation {
		openapi3.SchemaFormatValidationDisabled = true
	}

	if validate {
		if err = doc.Validate(openapi3.NewLoader().Context); err != nil {
			return "", nil, err
		}
	}

	if jsonContent, err = doc.MarshalJSON(); err != nil {
		return "", nil, err
	}

	if isFileYaml(uri) {
		yamlContent, err := yaml.JSONToYAML(jsonContent)
		return path.Base(u.Path), yamlContent, err
	} else {
		return path.Base(u.Path), jsonContent, err
	}
}

func isFileYaml(name string) bool {
	if strings.Contains(name, ".yaml") || strings.Contains(name, ".yml") {
		return true
	}
	return false
}

func GenerateAPIProxyDefFromOAS(name string, oasDocName string, skipPolicy bool, addCORS bool, oasGoogleAcessTokenScopeLiteral string, oasGoogleIdTokenAudLiteral string, oasGoogleIdTokenAudRef string, oasTargetUrlRef string) (err error) {

	if doc == nil {
		return fmt.Errorf("Open API document not loaded")
	}

	//load security schemes
	loadSecurityRequirements(doc.Components.SecuritySchemes)

	apiproxy.SetDisplayName(name)
	if doc.Info != nil {
		if doc.Info.Description != "" {
			apiproxy.SetDescription(doc.Info.Description)
		}
	}

	apiproxy.SetCreatedAt()
	apiproxy.SetLastModifiedAt()
	apiproxy.SetConfigurationVersion()
	apiproxy.AddTargetEndpoint("default")
	apiproxy.AddProxyEndpoint("default")

	if !skipPolicy {
		apiproxy.AddResource(oasDocName)
		apiproxy.AddPolicy("Validate-" + name + "-Schema")
	}

	u, err := GetEndpoint(doc)
	if err != nil {
		return err
	}

	if u.Path == "" {
		return fmt.Errorf("OpenAPI url is missing a path. Don't use https://api.example.com, instead try https://api.example.com/basePath")
	}

	apiproxy.SetBasePath(u.Path)

	//set a dynamic target url
	if oasTargetUrlRef != "" {
		targets.AddStepToPreFlowRequest("Set-Target-1")
		apiproxy.AddPolicy("Set-Target-1")
		generateSetTarget = true
	}

	targets.NewTargetEndpoint(u.Scheme+"://"+u.Hostname(), oasGoogleAcessTokenScopeLiteral, oasGoogleIdTokenAudLiteral, oasGoogleIdTokenAudRef)

	proxies.NewProxyEndpoint(u.Path)

	//add any preflow security schemes
	if securityScheme := getSecurityRequirements(doc.Security); securityScheme.SchemeName != "" {
		if securityScheme.APIKeyPolicy.APIKeyPolicyEnabled {
			proxies.AddStepToPreFlowRequest("Verify-API-Key-" + securityScheme.SchemeName)
		} else if securityScheme.OAuthPolicy.OAuthPolicyEnabled {
			apiproxy.AddPolicy("OAuth-v20-1")
		}
	}

	if addCORS {
		proxies.AddStepToPreFlowRequest("Add-CORS")
		apiproxy.AddPolicy("Add-CORS")
	}

	if !skipPolicy {
		proxies.AddStepToPreFlowRequest("OpenAPI-Spec-Validation-1")
	}

	if err = GenerateFlows(doc.Paths); err != nil {
		return err
	}

	for _, securityScheme := range securitySchemesList.SecuritySchemes {
		if securityScheme.APIKeyPolicy.APIKeyPolicyEnabled {
			apiproxy.AddPolicy("Verify-API-Key-" + securityScheme.SchemeName)
		} else if securityScheme.OAuthPolicy.OAuthPolicyEnabled {
			apiproxy.AddPolicy("OAuth-v20-1")
		}
	}

	return nil
}

func GetEndpoint(doc *openapi3.T) (u *url.URL, err error) {
	if doc.Servers == nil {
		return nil, fmt.Errorf("at least one server must be present")
	}

	return url.Parse(doc.Servers[0].URL)
}

func GetHTTPMethod(pathItem *openapi3.PathItem, keyPath string) map[string]pathDetailDef {

	pathMap := make(map[string]pathDetailDef)
	alternateOperationId := strings.ReplaceAll(keyPath, "\\", "_")

	if pathItem.Get != nil {
		getPathDetail := pathDetailDef{}
		if pathItem.Get.OperationID != "" {
			getPathDetail.OperationID = pathItem.Get.OperationID
		} else {
			getPathDetail.OperationID = "get_" + alternateOperationId
		}
		if pathItem.Get.Description != "" {
			getPathDetail.Description = pathItem.Get.Description
		}
		if pathItem.Get.Security != nil {
			securityRequirements := []openapi3.SecurityRequirement(*pathItem.Get.Security)
			getPathDetail.SecurityScheme = getSecurityRequirements(securityRequirements)
		}
		pathMap["get"] = getPathDetail
	}

	if pathItem.Post != nil {
		postPathDetail := pathDetailDef{}
		if pathItem.Post.OperationID != "" {
			postPathDetail.OperationID = pathItem.Post.OperationID
		} else {
			postPathDetail.OperationID = "post_" + alternateOperationId
		}
		if pathItem.Post.Description != "" {
			postPathDetail.Description = pathItem.Post.Description
		}
		if pathItem.Post.Security != nil {
			securityRequirements := []openapi3.SecurityRequirement(*pathItem.Post.Security)
			postPathDetail.SecurityScheme = getSecurityRequirements(securityRequirements)
		}
		pathMap["post"] = postPathDetail
	}

	if pathItem.Put != nil {
		putPathDetail := pathDetailDef{}
		if pathItem.Put.OperationID != "" {
			putPathDetail.OperationID = pathItem.Put.OperationID
		} else {
			putPathDetail.OperationID = "put_" + alternateOperationId
		}
		if pathItem.Put.Description != "" {
			putPathDetail.Description = pathItem.Put.Description
		}
		if pathItem.Put.Security != nil {
			securityRequirements := []openapi3.SecurityRequirement(*pathItem.Put.Security)
			putPathDetail.SecurityScheme = getSecurityRequirements(securityRequirements)
		}
		pathMap["put"] = putPathDetail
	}

	if pathItem.Patch != nil {
		patchPathDetail := pathDetailDef{}
		if pathItem.Patch.OperationID != "" {
			patchPathDetail.OperationID = pathItem.Patch.OperationID
		} else {
			patchPathDetail.OperationID = "patch_" + alternateOperationId
		}
		if pathItem.Patch.Description != "" {
			patchPathDetail.Description = pathItem.Patch.Description
		}
		if pathItem.Patch.Security != nil {
			securityRequirements := []openapi3.SecurityRequirement(*pathItem.Patch.Security)
			patchPathDetail.SecurityScheme = getSecurityRequirements(securityRequirements)
		}
		pathMap["patch"] = patchPathDetail
	}

	if pathItem.Delete != nil {
		deletePathDetail := pathDetailDef{}
		if pathItem.Delete.OperationID != "" {
			deletePathDetail.OperationID = pathItem.Delete.OperationID
		} else {
			deletePathDetail.OperationID = "delete_" + alternateOperationId
		}
		if pathItem.Delete.Description != "" {
			deletePathDetail.Description = pathItem.Delete.Description
		}
		if pathItem.Delete.Security != nil {
			securityRequirements := []openapi3.SecurityRequirement(*pathItem.Delete.Security)
			deletePathDetail.SecurityScheme = getSecurityRequirements(securityRequirements)
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

	if pathItem.Trace != nil {
		tracePathDetail := pathDetailDef{}
		if pathItem.Trace.OperationID != "" {
			tracePathDetail.OperationID = pathItem.Trace.OperationID
		} else {
			tracePathDetail.OperationID = "trace_" + alternateOperationId
		}
		if pathItem.Trace.Description != "" {
			tracePathDetail.Description = pathItem.Trace.Description
		}
		pathMap["trace"] = tracePathDetail
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

	return pathMap
}

func GenerateFlows(paths openapi3.Paths) (err error) {
	for keyPath := range paths {
		pathMap := GetHTTPMethod(paths[keyPath], keyPath)
		for method, pathDetail := range pathMap {
			proxies.AddFlow(pathDetail.OperationID, replacePathWithWildCard(keyPath), method, pathDetail.Description)
			if pathDetail.SecurityScheme.OAuthPolicy.OAuthPolicyEnabled {
				if err = proxies.AddStepToFlowRequest("OAuth-v20-1", pathDetail.OperationID); err != nil {
					return err
				}
			} else if pathDetail.SecurityScheme.APIKeyPolicy.APIKeyPolicyEnabled {
				if err = proxies.AddStepToFlowRequest("Verify-API-Key-"+pathDetail.SecurityScheme.SchemeName, pathDetail.OperationID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func GenerateSetTargetPolicy() bool {
	return generateSetTarget
}

func replacePathWithWildCard(keyPath string) string {
	re := regexp.MustCompile(`{(.*?)}`)
	if strings.ContainsAny(keyPath, "{") {
		return re.ReplaceAllLiteralString(keyPath, "*")
	} else {
		return keyPath
	}
}

func loadSecurityType(secSchemeName string, securityScheme openapi3.SecuritySchemeRef) (secScheme securitySchemesDef) {
	secScheme = securitySchemesDef{}
	apiKeyPolicy := apiKeyPolicyDef{}
	oAuthPolicy := oAuthPolicyDef{}

	if securityScheme.Value.Type == "oauth2" {
		secScheme.SchemeName = secSchemeName
		oAuthPolicy.OAuthPolicyEnabled = true
		apiKeyPolicy.APIKeyPolicyEnabled = false
		secScheme.OAuthPolicy = oAuthPolicy
	} else if securityScheme.Value.Type == "apiKey" {
		secScheme.SchemeName = secSchemeName
		apiKeyPolicy.APIKeyPolicyEnabled = true
		oAuthPolicy.OAuthPolicyEnabled = false
		apiKeyPolicy.APIKeyLocation = securityScheme.Value.In
		apiKeyPolicy.APIKeyName = securityScheme.Value.Name
		secScheme.APIKeyPolicy = apiKeyPolicy
	}

	return secScheme
}

func getSecurityType(secName string) securitySchemesDef {
	for _, securityScheme := range securitySchemesList.SecuritySchemes {
		if securityScheme.SchemeName == secName {
			return securityScheme
		}
	}
	return securitySchemesDef{}
}

func getSecurityRequirements(securityRequirements []openapi3.SecurityRequirement) securitySchemesDef {
	for _, secReq := range securityRequirements {
		for secReqName := range secReq {
			return getSecurityType(secReqName)
		}
	}
	return securitySchemesDef{}
}

func loadSecurityRequirements(securitySchemes openapi3.SecuritySchemes) {
	for secSchemeName, secScheme := range securitySchemes {
		securitySchemesList.SecuritySchemes = append(securitySchemesList.SecuritySchemes, loadSecurityType(secSchemeName, *secScheme))
	}
}

func GetSecuritySchemesList() []securitySchemesDef {
	return securitySchemesList.SecuritySchemes
}

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
	target "github.com/srinandan/apigeecli/bundlegen/targetendpoint"
)

type pathDetailDef struct {
	OperationID string
	Description string
	OAuthPolicy bool
	APIKeyPoicy bool
}

var generateOAuthPolicy, generateAPIKeyPolicy bool

var doc *openapi3.T

func LoadDocumentFromFile(filePath string, validate bool) (string, []byte, error) {
	var err error
	var jsonContent []byte

	doc, err = openapi3.NewLoader().LoadFromFile(filePath)
	if err != nil {
		return "", nil, err
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

func LoadDocumentFromURI(uri string, validate bool) (string, []byte, error) {
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

func GenerateAPIProxyDefFromOAS(name string, oasDocName string, skipPolicy bool, addCORS bool) (err error) {

	if doc == nil {
		return fmt.Errorf("Open API document not loaded")
	}

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

	apiproxy.SetBasePath(u.Path)

	target.NewTargetEndpoint(u.Scheme + "://" + u.Hostname())

	proxies.NewProxyEndpoint(u.Path)

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

	if GenerateAPIKeyPolicy() {
		apiproxy.AddPolicy("Verify-API-Key-1")
	}

	if GenerateOAuthPolicy() {
		apiproxy.AddPolicy("OAuth-v20-1")
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
			getPathDetail.OAuthPolicy, getPathDetail.APIKeyPoicy = getSecurityRequirements(securityRequirements)
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
			postPathDetail.OAuthPolicy, postPathDetail.APIKeyPoicy = getSecurityRequirements(securityRequirements)
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
			putPathDetail.OAuthPolicy, putPathDetail.APIKeyPoicy = getSecurityRequirements(securityRequirements)
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
			patchPathDetail.OAuthPolicy, patchPathDetail.APIKeyPoicy = getSecurityRequirements(securityRequirements)
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
			deletePathDetail.OAuthPolicy, deletePathDetail.APIKeyPoicy = getSecurityRequirements(securityRequirements)
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
			if pathDetail.OAuthPolicy {
				if err = proxies.AddStepToFlowRequest("OAuth-v20-1", pathDetail.OperationID); err != nil {
					return err
				}
			}
			if pathDetail.APIKeyPoicy {
				if err = proxies.AddStepToFlowRequest("Verify-API-Key-1", pathDetail.OperationID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func GenerateOAuthPolicy() bool {
	return generateOAuthPolicy
}

func GenerateAPIKeyPolicy() bool {
	return generateAPIKeyPolicy
}

func replacePathWithWildCard(keyPath string) string {
	re := regexp.MustCompile(`{(.*?)}`)
	if strings.ContainsAny(keyPath, "{") {
		return re.ReplaceAllLiteralString(keyPath, "*")
	} else {
		return keyPath
	}
}

func getSecurityType(name string) (oauth bool, apikey bool) {
	for schemeName, securityScheme := range doc.Components.SecuritySchemes {
		if schemeName == name {
			if securityScheme.Value.Type == "oauth2" {
				generateOAuthPolicy = true
				return true, false
			}
			if securityScheme.Value.Type == "apiKey" {
				generateAPIKeyPolicy = true
				return false, true
			}
		}
	}
	return false, false
}

func getSecurityRequirements(securityRequirements []openapi3.SecurityRequirement) (oauth bool, apikey bool) {
	for _, secReq := range securityRequirements {
		for secReqName := range secReq {
			return getSecurityType(secReqName)
		}
	}
	return false, false
}

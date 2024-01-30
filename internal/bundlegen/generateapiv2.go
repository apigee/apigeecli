// Copyright 2023 Google LLC
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
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"internal/apiclient"
	apiproxy "internal/bundlegen/apiproxydef"
	"internal/bundlegen/policies"
	"internal/bundlegen/proxies"
	targets "internal/bundlegen/targets"
	"internal/clilog"

	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"gopkg.in/yaml.v3"
)

var docModel *libopenapi.DocumentModel[v3.Document]

func LoadDocument(specBasePath string, specBaseURL string, specName string, validate bool) (contents []byte, err error) {
	config := datamodel.DocumentConfiguration{}
	var specBytes []byte
	var errs []error
	references := []byte("$ref")

	if specBasePath != "" {
		config = datamodel.DocumentConfiguration{
			AllowFileReferences:     true,
			AllowRemoteReferences:   true,
			BasePath:                specBasePath,
			BundleInlineRefs:        true,
			ExtractRefsSequentially: true,
		}
		specBytes, err = os.ReadFile(filepath.Join(specBasePath, specName))
		if err != nil {
			return nil, err
		}
	} else {
		baseURL, _ := url.Parse(specBaseURL)
		config = datamodel.DocumentConfiguration{
			AllowFileReferences:     true,
			AllowRemoteReferences:   true,
			BaseURL:                 baseURL,
			BundleInlineRefs:        true,
			ExtractRefsSequentially: true,
		}
		specURL, _ := url.JoinPath(specBaseURL, specName)
		resp, err := apiclient.DownloadFile(specURL, false)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		specBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}
	document, err := libopenapi.NewDocumentWithConfiguration(specBytes, &config)
	if err != nil {
		return nil, err
	}
	if index := bytes.Index(specBytes, references); index != -1 && validate {
		clilog.Warning.Println("found references in the spec. disabling validation.")
		validate = false
	}
	if validate {
		docValidator, err := validator.NewValidator(document)
		if err != nil {
			return nil, fmt.Errorf("failed to create validator: %v", err)
		}
		valid, validationErrs := docValidator.ValidateDocument()
		if !valid {
			for i, e := range validationErrs {
				validationErr := fmt.Errorf("#%d - %d:%d, - type: %s, failure: %s, fix: %s", i, e.SpecLine,
					e.SpecCol, e.ValidationType, e.Message, e.HowToFix)
				errs = append(errs, validationErr)
			}
			return nil, errors.Join(errs...)
		}
	}
	docModel, errs = document.BuildV3Model()
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return specBytes, nil
}

func GetModelVersion() string {
	if docModel.Model.Version != "" {
		return docModel.Model.Version
	} else {
		return ""
	}
}

func GenerateAPIProxyDefFromOASv2(name string,
	basePath string,
	oasDocName string,
	skipPolicy bool,
	addCORS bool,
	integrationEndpoint bool,
	oasGoogleAcessTokenScopeLiteral string,
	oasGoogleIDTokenAudLiteral string,
	oasGoogleIDTokenAudRef string,
	oasTargetURLRef string,
	targetURL string,
) (err error) {

	if docModel == nil {
		return fmt.Errorf("the Open API document not loaded")
	}

	// load security schemes
	loadSecurityRequirementsv2(docModel.Model.Components.SecuritySchemes)

	apiproxy.SetDisplayName(name)
	if docModel.Model.Info != nil {
		if docModel.Model.Info.Description != "" {
			apiproxy.SetDescription(docModel.Model.Info.Description)
		}
	}

	apiproxy.SetCreatedAt()
	apiproxy.SetLastModifiedAt()
	apiproxy.SetConfigurationVersion()
	if integrationEndpoint {
		apiproxy.AddIntegrationEndpoint("default")
	} else {
		apiproxy.AddTargetEndpoint(NoAuthTargetName)
	}

	apiproxy.AddProxyEndpoint("default")

	if !skipPolicy {
		apiproxy.AddResource(oasDocName, "oas")
		apiproxy.AddPolicy("Validate-" + name + "-Schema")
	}

	u, err := getEndpointv2(docModel.Model.Servers)
	if err != nil {
		return err
	}

	if basePath == "" && u.Path == "" {
		return fmt.Errorf("the OpenAPI url is missing a path. Add basePath as a parameter or change the spec." +
			"Don't use https://api.example.com, instead try https://api.example.com/basePath")
	}

	if basePath != "" {
		apiproxy.SetBasePath(basePath)
	} else {
		apiproxy.SetBasePath(u.Path)
	}

	// decide on the type of target
	if integrationEndpoint { // assume an integration endpoint
		proxies.AddStepToPreFlowRequest("set-integration-request")
		apiproxy.AddPolicy("set-integration-request")
		proxies.NewProxyEndpoint(u.Path, false)
	} else {
		// if target is not set, derive it from the OAS file
		if targetURL == "" {
			targets.NewTargetEndpoint(NoAuthTargetName,
				u.Scheme+"://"+u.Hostname()+u.Path,
				oasGoogleAcessTokenScopeLiteral,
				oasGoogleIDTokenAudLiteral,
				oasGoogleIDTokenAudRef)
		} else { // an explicit target url is set
			if _, err = url.Parse(targetURL); err != nil {
				return fmt.Errorf("invalid target url: %v", err)
			}
			targets.NewTargetEndpoint(NoAuthTargetName,
				targetURL,
				oasGoogleAcessTokenScopeLiteral,
				oasGoogleIDTokenAudLiteral,
				oasGoogleIDTokenAudRef)
		}

		// set a dynamic target url
		if oasTargetURLRef != "" {
			targets.AddStepToPreFlowRequest("Set-Target-1", NoAuthTargetName)
			apiproxy.AddPolicy("Set-Target-1")
			generateSetTarget = true
		}

		if basePath != "" {
			proxies.NewProxyEndpoint(basePath, true)
		} else {
			proxies.NewProxyEndpoint(u.Path, true)
		}
	}

	// add any preflow security schemes
	if securityScheme := getSecurityRequirementsv2(docModel.Model.Security); securityScheme.SchemeName != "" {
		if securityScheme.APIKeyPolicy.APIKeyPolicyEnabled {
			proxies.AddStepToPreFlowRequest("Verify-API-Key-" + securityScheme.SchemeName)
		} else if securityScheme.OAuthPolicy.OAuthPolicyEnabled {
			proxies.AddStepToPreFlowRequest("OAuth-v20-1")
		}
	}

	// add any preflow quota or rate limit policies
	if docModel.Model.Extensions != nil {
		spikeArrestList, quotaList, err := processPreFlowExtensionsv2(docModel.Model.Extensions)
		if err != nil {
			return err
		}
		if len(spikeArrestList) > 0 {
			for _, spikeArrest := range spikeArrestList {
				proxies.AddStepToPreFlowRequest("Spike-Arrest-" + spikeArrest.SpikeArrestName)
			}
		}
		if len(quotaList) > 0 {
			for _, quota := range quotaList {
				proxies.AddStepToPreFlowRequest("Quota-" + quota.QuotaName)
			}
		}
	}

	if addCORS {
		proxies.AddStepToPreFlowRequest("Add-CORS")
		apiproxy.AddPolicy("Add-CORS")
	}

	if !skipPolicy {
		proxies.AddStepToPreFlowRequest("OpenAPI-Spec-Validation-1")
	}

	if err = generateFlowsv2(docModel.Model.Paths); err != nil {
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

func loadSecurityRequirementsv2(securitySchemes *orderedmap.Map[string, *v3.SecurityScheme]) (err error) {
	for first := securitySchemes.First(); first != nil; first = first.Next() {
		securityScheme := first.Value()
		securitySchemesList.SecuritySchemes = append(securitySchemesList.SecuritySchemes, loadSecurityTypev2(first.Key(), securityScheme))
	}
	return nil
}

func getEndpointv2(servers []*v3.Server) (u *url.URL, err error) {
	if servers == nil {
		return nil, fmt.Errorf("at least one server must be present")
	}
	return url.Parse(servers[0].URL)
}

func loadSecurityTypev2(secSchemeName string, securityScheme *v3.SecurityScheme) (secScheme securitySchemesDef) {
	secScheme = securitySchemesDef{}
	apiKeyPolicy := apiKeyPolicyDef{}
	oAuthPolicy := oAuthPolicyDef{}
	jwtPolicy := jwtPolicyDef{}

	// this policy does not apply for OAS 3, used only with Endpoints/Swagger
	jwtPolicy.JWTPolicyEnabled = false

	if securityScheme.Type == "oauth2" {
		secScheme.SchemeName = secSchemeName
		oAuthPolicy.OAuthPolicyEnabled = true
		apiKeyPolicy.APIKeyPolicyEnabled = false
		if securityScheme.Flows != nil {
			if securityScheme.Flows.Implicit != nil {
				oAuthPolicy.Scope = readScopesv2(securityScheme.Flows.Implicit.Scopes)
			} else if securityScheme.Flows.Password != nil {
				oAuthPolicy.Scope = readScopesv2(securityScheme.Flows.Password.Scopes)
			} else if securityScheme.Flows.ClientCredentials.Scopes != nil {
				oAuthPolicy.Scope = readScopesv2(securityScheme.Flows.ClientCredentials.Scopes)
			} else if securityScheme.Flows.AuthorizationCode.Scopes != nil {
				oAuthPolicy.Scope = readScopesv2(securityScheme.Flows.AuthorizationCode.Scopes)
			}
		}
		secScheme.OAuthPolicy = oAuthPolicy
	} else if securityScheme.Type == "apiKey" {
		secScheme.SchemeName = secSchemeName
		apiKeyPolicy.APIKeyPolicyEnabled = true
		oAuthPolicy.OAuthPolicyEnabled = false
		apiKeyPolicy.APIKeyLocation = securityScheme.In
		apiKeyPolicy.APIKeyName = securityScheme.Name
		secScheme.APIKeyPolicy = apiKeyPolicy
	}

	return secScheme
}

func readScopesv2(scopes *orderedmap.Map[string, string]) string {
	scopeString := ""
	for first := scopes.First(); first != nil; first = first.Next() {
		scopeString = first.Key() + " " + scopeString
	}
	return strings.TrimSpace(scopeString)
}

func getSecurityRequirementsv2(securityRequirements []*base.SecurityRequirement) securitySchemesDef {
	for _, secReq := range securityRequirements {
		for first := secReq.Requirements.First(); first != nil; first = first.Next() {
			return getSecurityType(first.Key())
		}
	}
	return securitySchemesDef{}
}

func processPreFlowExtensionsv2(extension *orderedmap.Map[string, *yaml.Node]) ([]spikeArrestDef, []quotaDef, error) {
	var err error
	spikeArrestList := []spikeArrestDef{}
	quotaList := []quotaDef{}

	for first := extension.First(); first != nil; first = first.Next() {
		extensionName := first.Key()
		extensionValue := first.Value()
		if extensionName == "x-google-ratelimit" {
			// process ratelimit
			spikeArrests, err := getSpikeArrestDefinitionv2(extensionValue)
			if err != nil {
				return []spikeArrestDef{}, []quotaDef{}, err
			}
			spikeArrestList = append(spikeArrestList, spikeArrests...)
		}
		if extensionName == "x-google-quota" {
			// process quota
			quotas, err := getQuotaDefinitionv2(extensionValue)
			if err != nil {
				return []spikeArrestDef{}, []quotaDef{}, err
			}
			quotaList = append(quotaList, quotas...)
		}
	}

	return spikeArrestList, quotaList, err
}

func getSpikeArrestDefinitionv2(extension *yaml.Node) ([]spikeArrestDef, error) {
	var jsonArrayMap []map[string]interface{}
	spikeArrests := []spikeArrestDef{}

	err := extension.Decode(&jsonArrayMap)
	if err != nil {
		return spikeArrests, err
	}

	for _, m := range jsonArrayMap {
		spikeArrest := spikeArrestDef{}
		for k, v := range m {
			if k == "identifier-ref" {
				spikeArrest.SpikeArrestIdentifierRef = fmt.Sprintf("%v", v)
			} else if k == "name" {
				spikeArrest.SpikeArrestName = fmt.Sprintf("%v", v)
				spikeArrest.SpikeArrestEnabled = true
			} else if k == "rate-literal" {
				spikeArrest.SpikeArrestRateLiteral = fmt.Sprintf("%v", v)
			} else if k == "rate-ref" {
				spikeArrest.SpikeArrestRateRef = fmt.Sprintf("%v", v)
			}
		}
		spikeArrests = append(spikeArrests, spikeArrest)
	}

	for _, s := range spikeArrests {
		if s.SpikeArrestName == "" {
			return nil, fmt.Errorf("x-google-ratelimit extension must have a name")
		}
		if s.SpikeArrestIdentifierRef == "" {
			return nil, fmt.Errorf("x-google-ratelimit extension must have an identifier-ref")
		}
		if s.SpikeArrestRateLiteral == "" && s.SpikeArrestRateRef == "" {
			return nil, fmt.Errorf("x-google-ratelimit extension must have either rate-ref or rate-literal")
		}
		// store policy XML contents
		spikeArrestPolicyContent[s.SpikeArrestName] = policies.AddSpikeArrestPolicy("Spike-Arrest-"+s.SpikeArrestName,
			s.SpikeArrestIdentifierRef,
			s.SpikeArrestRateRef,
			s.SpikeArrestRateLiteral)
	}

	return spikeArrests, nil
}

func getQuotaDefinitionv2(extension *yaml.Node) ([]quotaDef, error) {
	var jsonArrayMap []map[string]interface{}
	quotas := []quotaDef{}

	err := extension.Decode(&jsonArrayMap)
	if err != nil {
		return quotas, err
	}

	for _, m := range jsonArrayMap {
		quota := quotaDef{}
		for k, v := range m {
			if k == "name" {
				quota.QuotaName = fmt.Sprintf("%v", v)
				quota.QuotaEnabled = true
			} else if k == "useQuotaConfigInAPIProduct" {
				quota.QuotaConfigStepName = fmt.Sprintf("%v", v)
			} else if k == "allow-ref" {
				quota.QuotaAllowRef = fmt.Sprintf("%v", v)
			} else if k == "allow-literal" {
				quota.QuotaAllowLiteral = fmt.Sprintf("%v", v)
			} else if k == "interval-ref" {
				quota.QuotaIntervalRef = fmt.Sprintf("%v", v)
			} else if k == "interval-literal" {
				quota.QuotaIntervalLiteral = fmt.Sprintf("%v", v)
			} else if k == "timeunit-ref" {
				quota.QuotaTimeUnitRef = fmt.Sprintf("%v", v)
			} else if k == "timeunit-literal" {
				quota.QuotaTimeUnitLiteral = fmt.Sprintf("%v", v)
			} else if k == "identifier-ref" {
				quota.QuotaIdentifierRef = fmt.Sprintf("%v", v)
			} else if k == "identifier-literal" {
				quota.QuotaIdentiferLiteral = fmt.Sprintf("%v", v)
			}
		}
		quotas = append(quotas, quota)
	}

	for _, q := range quotas {
		if q.QuotaName == "" {
			return nil, fmt.Errorf("x-google-quota extension must have a name")
		}
		if q.QuotaConfigStepName == "" {
			if q.QuotaAllowLiteral == "" && q.QuotaAllowRef == "" {
				return nil, fmt.Errorf("x-google-quota extension must have either allow-ref or allow-literal")
			}
			if q.QuotaIntervalLiteral == "" && q.QuotaIntervalRef == "" {
				return nil, fmt.Errorf("x-google-quota extension must have either interval-ref or interval-literal")
			}
			if q.QuotaTimeUnitLiteral == "" && q.QuotaTimeUnitRef == "" {
				return nil, fmt.Errorf("x-google-quota extension must have either timeunit-ref or timeunit-literal")
			}
		}
		// store policy XML contents
		quotaPolicyContent[q.QuotaName] = policies.AddQuotaPolicy(
			"Quota-"+q.QuotaName,
			q.QuotaConfigStepName,
			q.QuotaAllowRef,
			q.QuotaAllowLiteral,
			q.QuotaIntervalRef,
			q.QuotaIntervalLiteral,
			q.QuotaTimeUnitRef,
			q.QuotaTimeUnitLiteral,
			q.QuotaIdentifierRef,
			q.QuotaIdentiferLiteral)
	}

	return quotas, nil
}

func generateFlowsv2(paths *v3.Paths) (err error) {
	for first := paths.PathItems.First(); first != nil; first = first.Next() {
		pathMap, err := getHTTPMethodv2(first.Value(), first.Key())
		if err != nil {
			return err
		}
		for method, pathDetail := range pathMap {
			proxies.AddFlow(pathDetail.OperationID, replacePathWithWildCard(first.Key()), method, pathDetail.Description)
			if pathDetail.SecurityScheme.OAuthPolicy.OAuthPolicyEnabled {
				if err = proxies.AddStepToFlowRequest("OAuth-v20-1", pathDetail.OperationID); err != nil {
					return err
				}
			} else if pathDetail.SecurityScheme.APIKeyPolicy.APIKeyPolicyEnabled {
				if err = proxies.AddStepToFlowRequest("Verify-API-Key-"+pathDetail.SecurityScheme.SchemeName, pathDetail.OperationID); err != nil {
					return err
				}
			}
			for _, s := range pathDetail.SpikeArrest {
				if s.SpikeArrestEnabled {
					if err = proxies.AddStepToFlowRequest("Spike-Arrest-"+s.SpikeArrestName, pathDetail.OperationID); err != nil {
						return err
					}
				}
			}
			for _, q := range pathDetail.Quota {
				if q.QuotaEnabled {
					if err = proxies.AddStepToFlowRequest("Quota-"+q.QuotaName, pathDetail.OperationID); err != nil {
						return err
					}
				}
			}
		}
	}

	return err
}

func getHTTPMethodv2(pathItem *v3.PathItem, keyPath string) (map[string]pathDetailDef, error) {
	var err error
	pathMap := make(map[string]pathDetailDef)
	alternateOperationId := strings.ReplaceAll(keyPath, "\\", "_")

	if pathItem.Get != nil {
		getPathDetail := pathDetailDef{}
		if pathItem.Get.OperationId != "" {
			getPathDetail.OperationID = pathItem.Get.OperationId
		} else {
			getPathDetail.OperationID = "get_" + alternateOperationId
		}
		if pathItem.Get.Description != "" {
			getPathDetail.Description = pathItem.Get.Description
		}
		if pathItem.Get.Security != nil {
			getPathDetail.SecurityScheme = getSecurityRequirementsv2(pathItem.Get.Security)
		}
		// check for google extensions
		if pathItem.Get.Extensions != nil {
			if getPathDetail, err = processPathExtensionsv2(pathItem.Get.Extensions, getPathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["get"] = getPathDetail
	}

	if pathItem.Post != nil {
		postPathDetail := pathDetailDef{}
		if pathItem.Post.OperationId != "" {
			postPathDetail.OperationID = pathItem.Post.OperationId
		} else {
			postPathDetail.OperationID = "post_" + alternateOperationId
		}
		if pathItem.Post.Description != "" {
			postPathDetail.Description = pathItem.Post.Description
		}
		if pathItem.Post.Security != nil {
			postPathDetail.SecurityScheme = getSecurityRequirementsv2(pathItem.Post.Security)
		}
		// check for google extensions
		if pathItem.Post.Extensions != nil {
			if postPathDetail, err = processPathExtensionsv2(pathItem.Post.Extensions, postPathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["post"] = postPathDetail
	}

	if pathItem.Put != nil {
		putPathDetail := pathDetailDef{}
		if pathItem.Put.OperationId != "" {
			putPathDetail.OperationID = pathItem.Put.OperationId
		} else {
			putPathDetail.OperationID = "put_" + alternateOperationId
		}
		if pathItem.Put.Description != "" {
			putPathDetail.Description = pathItem.Put.Description
		}
		if pathItem.Put.Security != nil {
			putPathDetail.SecurityScheme = getSecurityRequirementsv2(pathItem.Put.Security)
		}
		// check for google extensions
		if pathItem.Put.Extensions != nil {
			if putPathDetail, err = processPathExtensionsv2(pathItem.Put.Extensions, putPathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["put"] = putPathDetail
	}

	if pathItem.Patch != nil {
		patchPathDetail := pathDetailDef{}
		if pathItem.Patch.OperationId != "" {
			patchPathDetail.OperationID = pathItem.Patch.OperationId
		} else {
			patchPathDetail.OperationID = "patch_" + alternateOperationId
		}
		if pathItem.Patch.Description != "" {
			patchPathDetail.Description = pathItem.Patch.Description
		}
		if pathItem.Patch.Security != nil {
			patchPathDetail.SecurityScheme = getSecurityRequirementsv2(pathItem.Patch.Security)
		}
		// check for google extensions
		if pathItem.Patch.Extensions != nil {
			if patchPathDetail, err = processPathExtensionsv2(pathItem.Patch.Extensions, patchPathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["patch"] = patchPathDetail
	}

	if pathItem.Delete != nil {
		deletePathDetail := pathDetailDef{}
		if pathItem.Delete.OperationId != "" {
			deletePathDetail.OperationID = pathItem.Delete.OperationId
		} else {
			deletePathDetail.OperationID = "delete_" + alternateOperationId
		}
		if pathItem.Delete.Description != "" {
			deletePathDetail.Description = pathItem.Delete.Description
		}
		if pathItem.Delete.Security != nil {
			deletePathDetail.SecurityScheme = getSecurityRequirementsv2(pathItem.Delete.Security)
		}
		// check for google extensions
		if pathItem.Delete.Extensions != nil {
			if deletePathDetail, err = processPathExtensionsv2(pathItem.Delete.Extensions, deletePathDetail); err != nil {
				return nil, err
			}
		}
		pathMap["delete"] = deletePathDetail
	}

	if pathItem.Options != nil {
		optionsPathDetail := pathDetailDef{}
		if pathItem.Options.OperationId != "" {
			optionsPathDetail.OperationID = pathItem.Options.OperationId
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
		if pathItem.Trace.OperationId != "" {
			tracePathDetail.OperationID = pathItem.Trace.OperationId
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
		if pathItem.Head.OperationId != "" {
			headPathDetail.OperationID = pathItem.Head.OperationId
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

func processPathExtensionsv2(extensions *orderedmap.Map[string, *yaml.Node], pathDetail pathDetailDef) (pathDetailDef, error) {
	var errs []error
	var err error
	for first := extensions.First(); first != nil; first = first.Next() {
		if first.Key() == "x-google-ratelimit" {
			// process rate limit
			pathDetail.SpikeArrest, err = getSpikeArrestDefinitionv2(first.Value())
		}
		if first.Key() == "x-google-quota" {
			// process quota
			pathDetail.Quota, err = getQuotaDefinitionv2(first.Value())
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return pathDetail, errors.Join(errs...)
}

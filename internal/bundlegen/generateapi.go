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
	"strconv"
	"strings"

	"internal/clilog"

	apiproxy "internal/bundlegen/apiproxydef"
	"internal/bundlegen/policies"
	"internal/bundlegen/proxies"
	targets "internal/bundlegen/targets"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
)

type pathDetailDef struct {
	OperationID    string
	Description    string
	Path           string
	Verb           string
	AssignMessage  string
	Backend        backendDef
	SecurityScheme securitySchemesDef
	SpikeArrest    []spikeArrestDef
	Quota          []quotaDef
}

type spikeArrestDef struct {
	SpikeArrestEnabled       bool
	SpikeArrestName          string
	SpikeArrestIdentifierRef string
	SpikeArrestRateRef       string
	SpikeArrestRateLiteral   string
}

type quotaDef struct {
	QuotaEnabled          bool
	QuotaName             string
	QuotaAllowRef         string
	QuotaAllowLiteral     string
	QuotaIntervalRef      string
	QuotaIntervalLiteral  string
	QuotaTimeUnitRef      string
	QuotaTimeUnitLiteral  string
	QuotaConfigStepName   string
	QuotaIdentifierRef    string
	QuotaIdentiferLiteral string
}

type oAuthPolicyDef struct {
	OAuthPolicyEnabled bool
	Scope              string
}

type securitySchemesListDef struct {
	SecuritySchemes []securitySchemesDef
}

type securitySchemesDef struct {
	SchemeName   string
	OAuthPolicy  oAuthPolicyDef
	APIKeyPolicy apiKeyPolicyDef
	JWTPolicy    jwtPolicyDef
}

type apiKeyPolicyDef struct {
	APIKeyPolicyEnabled bool
	APIKeyLocation      string
	APIKeyName          string
}

type jwtPolicyDef struct {
	JWTPolicyEnabled bool
	JwkUri           string
	Issuer           string
	Audience         string
	Source           string
	Location         map[string]string // only one location supported for now
}

var generateSetTarget bool

var securitySchemesList = securitySchemesListDef{}

var (
	quotaPolicyContent       = map[string]string{}
	spikeArrestPolicyContent = map[string]string{}
)

var doc *openapi3.T

func LoadDocumentFromFile(filePath string, validate bool, formatValidation bool) (string, []byte, error) {
	var err error
	var jsonContent []byte

	// see ./test/circular-reference.json and https://github.com/apigee/apigeecli/issues/199
	openapi3.CircularReferenceCounter = 20

	doc, err = openapi3.NewLoader().LoadFromFile(filePath)
	if err != nil {
		clilog.Error.Println(err)
		return "", nil, err
	}

	// add custom string definitions
	openapi3.DefineStringFormat("uuid", openapi3.FormatOfStringForUUIDOfRFC4122)

	if formatValidation {
		openapi3.EnableSchemaFormatValidation()
	} else {
		openapi3.DisableSchemaFormatValidation()
	}

	if validate {
		if err = doc.Validate(openapi3.NewLoader().Context); err != nil {
			clilog.Error.Println(err)
			return "", nil, err
		}
	}

	if jsonContent, err = doc.MarshalJSON(); err != nil {
		clilog.Error.Println(err)
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
		clilog.Error.Println(err)
		return "", nil, err
	}

	// see ./test/circular-reference.json and https://github.com/apigee/apigeecli/issues/199
	openapi3.CircularReferenceCounter = 20

	doc, err = openapi3.NewLoader().LoadFromURI(u)
	if err != nil {
		return "", nil, err
	}

	// add custom string definitions
	openapi3.DefineStringFormat("uuid", openapi3.FormatOfStringForUUIDOfRFC4122)

	if !formatValidation {
		openapi3.EnableSchemaFormatValidation()
	}

	if validate {
		if err = doc.Validate(openapi3.NewLoader().Context); err != nil {
			clilog.Error.Println(err)
			return "", nil, err
		}
	}

	if jsonContent, err = doc.MarshalJSON(); err != nil {
		clilog.Error.Println(err)
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
	if filepath.Ext(name) == ".yaml" || filepath.Ext(name) == ".yml" {
		return true
	}
	return false
}

func GenerateAPIProxyDefFromOAS(name string,
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
	if doc == nil {
		return fmt.Errorf("the Open API document not loaded")
	}

	// load security schemes
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

	u, err := getEndpoint(doc)
	if err != nil {
		return err
	}

	if u.Path == "" {
		return fmt.Errorf("the OpenAPI url is missing a path. Don't use https://api.example.com, instead try https://api.example.com/basePath")
	}

	apiproxy.SetBasePath(u.Path)

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

		proxies.NewProxyEndpoint(u.Path, true)
	}

	// add any preflow security schemes
	if securityScheme := getSecurityRequirements(doc.Security); securityScheme.SchemeName != "" {
		if securityScheme.APIKeyPolicy.APIKeyPolicyEnabled {
			proxies.AddStepToPreFlowRequest("Verify-API-Key-" + securityScheme.SchemeName)
		} else if securityScheme.OAuthPolicy.OAuthPolicyEnabled {
			proxies.AddStepToPreFlowRequest("OAuth-v20-1")
		}
	}

	// add any preflow quota or rate limit policies
	if doc.Extensions != nil {
		spikeArrestList, quotaList, err := processPreFlowExtensions(doc.Extensions)
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

	if err = generateFlows(doc.Paths); err != nil {
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

func getEndpoint(doc *openapi3.T) (u *url.URL, err error) {
	if doc.Servers == nil {
		return nil, fmt.Errorf("at least one server must be present")
	}

	return url.Parse(doc.Servers[0].URL)
}

func getHTTPMethod(pathItem *openapi3.PathItem, keyPath string) (map[string]pathDetailDef, error) {
	var err error
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
		// check for google extensions
		if pathItem.Get.Extensions != nil {
			if getPathDetail, err = processPathExtensions(pathItem.Get.Extensions, getPathDetail); err != nil {
				return nil, err
			}
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
		// check for google extensions
		if pathItem.Post.Extensions != nil {
			if postPathDetail, err = processPathExtensions(pathItem.Post.Extensions, postPathDetail); err != nil {
				return nil, err
			}
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
		// check for google extensions
		if pathItem.Put.Extensions != nil {
			if putPathDetail, err = processPathExtensions(pathItem.Put.Extensions, putPathDetail); err != nil {
				return nil, err
			}
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
		// check for google extensions
		if pathItem.Patch.Extensions != nil {
			if patchPathDetail, err = processPathExtensions(pathItem.Patch.Extensions, patchPathDetail); err != nil {
				return nil, err
			}
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
		// check for google extensions
		if pathItem.Delete.Extensions != nil {
			if deletePathDetail, err = processPathExtensions(pathItem.Delete.Extensions, deletePathDetail); err != nil {
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

	return pathMap, nil
}

func generateFlows(paths openapi3.Paths) (err error) {
	for keyPath := range paths {
		pathMap, err := getHTTPMethod(paths[keyPath], keyPath)
		if err != nil {
			return err
		}
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
	jwtPolicy := jwtPolicyDef{}

	// this policy does not apply for OAS 3, used only with Endpoints/Swagger
	jwtPolicy.JWTPolicyEnabled = false

	if securityScheme.Value.Type == "oauth2" {
		secScheme.SchemeName = secSchemeName
		oAuthPolicy.OAuthPolicyEnabled = true
		apiKeyPolicy.APIKeyPolicyEnabled = false
		if securityScheme.Value.Flows != nil {
			if securityScheme.Value.Flows.Implicit != nil {
				oAuthPolicy.Scope = readScopes(securityScheme.Value.Flows.Implicit.Scopes)
			} else if securityScheme.Value.Flows.Password != nil {
				oAuthPolicy.Scope = readScopes(securityScheme.Value.Flows.Password.Scopes)
			} else if securityScheme.Value.Flows.ClientCredentials.Scopes != nil {
				oAuthPolicy.Scope = readScopes(securityScheme.Value.Flows.ClientCredentials.Scopes)
			} else if securityScheme.Value.Flows.AuthorizationCode.Scopes != nil {
				oAuthPolicy.Scope = readScopes(securityScheme.Value.Flows.AuthorizationCode.Scopes)
			}
		}
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

func getQuotaDefinition(i interface{}) ([]quotaDef, error) {
	var jsonArrayMap []map[string]interface{}

	quotas := []quotaDef{}
	str := fmt.Sprintf("%s", i)
	jsonArrayMap = parseExtension(str)

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
			if q.QuotaIdentiferLiteral == "" && q.QuotaIdentifierRef == "" {
				return nil, fmt.Errorf("x-google-quota extension must have either identifier-ref or identifier-literal")
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

func getSpikeArrestDefinition(i interface{}) ([]spikeArrestDef, error) {
	var jsonArrayMap []map[string]interface{}

	spikeArrests := []spikeArrestDef{}
	str := fmt.Sprintf("%s", i)
	jsonArrayMap = parseExtension(str)

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

func processPathExtensions(extensions map[string]interface{}, pathDetail pathDetailDef) (pathDetailDef, error) {
	var err error
	for extensionName, extensionValue := range extensions {
		if extensionName == "x-google-ratelimit" {
			// process ratelimit
			pathDetail.SpikeArrest, err = getSpikeArrestDefinition(extensionValue)
		}
		if extensionName == "x-google-quota" {
			// process quota
			pathDetail.Quota, err = getQuotaDefinition(extensionValue)
		}
	}
	return pathDetail, err
}

func processPreFlowExtensions(extensions map[string]interface{}) ([]spikeArrestDef, []quotaDef, error) {
	var err error
	spikeArrestList := []spikeArrestDef{}
	quotaList := []quotaDef{}

	for extensionName, extensionValue := range extensions {
		if extensionName == "x-google-ratelimit" {
			// process ratelimit
			spikeArrests, err := getSpikeArrestDefinition(extensionValue)
			if err != nil {
				return []spikeArrestDef{}, []quotaDef{}, err
			}
			spikeArrestList = append(spikeArrestList, spikeArrests...)
		}
		if extensionName == "x-google-quota" {
			// process quota
			quotas, err := getQuotaDefinition(extensionValue)
			if err != nil {
				return []spikeArrestDef{}, []quotaDef{}, err
			}
			quotaList = append(quotaList, quotas...)
		}
	}
	return spikeArrestList, quotaList, err
}

func GetSpikeArrestPolicies() map[string]string {
	return spikeArrestPolicyContent
}

func GetQuotaPolicies() map[string]string {
	return quotaPolicyContent
}

func readScopes(scopes map[string]string) string {
	scopeString := ""
	for scopeName := range scopes {
		scopeString = scopeName + " " + scopeString
	}
	return strings.TrimSpace(scopeString)
}

func parseExtension(extension string) []map[string]interface{} {
	var result []map[string]interface{}
	pattern := `map\[(.*?)\]`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(extension, -1)
	for _, match := range matches {
		kvPairs := parseKeyValuePairs(match[1])
		result = append(result, kvPairs)
	}
	return result
}

func parseKeyValuePairs(match string) map[string]interface{} {
	// HACK: numbers in yaml are represented as %!s(float64=xx)
	match = strings.ReplaceAll(match, "%!s(float64=", "")
	match = strings.ReplaceAll(match, ")", "")

	keyValuePairs := make(map[string]interface{})
	keyValuePattern := `\b([A-Za-z0-9\-_]+)\s*:\s*([A-Za-z0-9\-_\.]+)\b`
	keyValueRe := regexp.MustCompile(keyValuePattern)
	kvMatches := keyValueRe.FindAllStringSubmatch(match, -1)
	for _, kvMatch := range kvMatches {
		key := kvMatch[1]
		value := kvMatch[2]
		if intValue, err := strconv.Atoi(value); err == nil {
			keyValuePairs[key] = intValue
		} else {
			keyValuePairs[key] = value
		}
	}
	return keyValuePairs
}

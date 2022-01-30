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

package policies

import (
	"regexp"
	"strings"
)

var oasPolicyTemplate = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<OASValidation continueOnError="false" enabled="true" name="OpenAPI-Spec-Validation-1">
    <DisplayName>OpenAPI Spec Validation-1</DisplayName>
    <Properties/>
    <Source>request</Source>
    <OASResource>oas://{PolicyName}</OASResource>
</OASValidation>`

var verifyApiKeyPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<VerifyAPIKey async="false" continueOnError="false" enabled="true" name="Verify-API-Key-1">
    <DisplayName>Verify-API-Key-1</DisplayName>
    <Properties/>
    <APIKey ref="request.queryparam.apikey"/>
</VerifyAPIKey>`

var oauth2Policy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<OAuthV2 async="false" continueOnError="false" enabled="true" name="OAuth-v20-1">
    <DisplayName>OAuth v2.0-1</DisplayName>
    <Properties/>
    <Attributes/>
    <ExternalAuthorization>false</ExternalAuthorization>
    <Operation>VerifyAccessToken</Operation>
    <SupportedGrantTypes/>
    <GenerateResponse enabled="true"/>
    <Tokens/>
</OAuthV2>`

var corsPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<CORS async="false" continueOnError="false" enabled="true" name="Add-CORS">
    <DisplayName>Add CORS</DisplayName>
    <AllowOrigins>{request.header.origin}</AllowOrigins>
    <AllowMethods>GET, PUT, POST, DELETE</AllowMethods>
    <AllowHeaders>origin, x-requested-with, accept, content-type</AllowHeaders>
    <ExposeHeaders>*</ExposeHeaders>
    <MaxAge>3628800</MaxAge>
    <AllowCredentials>false</AllowCredentials>
    <GeneratePreflightResponse>true</GeneratePreflightResponse>
    <IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
</CORS>`

var setTargetEndpointPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<AssignMessage async="false" continueOnError="false" enabled="true" name="Set-Target-1">
    <AssignVariable>
        <Name>target.url</Name>
        <Ref>dynamic.target.url</Ref>
    </AssignVariable>
    <IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
</AssignMessage>`

func AddOpenAPIValidatePolicy(name string) string {
	return replaceTemplateWithPolicy(name)
}

func AddVerifyApiKeyPolicy(location string, policyName string, keyName string) string {
	var apiKeyLocation string
	if location == "query" {
		apiKeyLocation = "request.queryparam." + keyName
	} else {
		apiKeyLocation = "request.header." + keyName
	}
	tmp := strings.Replace(verifyApiKeyPolicy, "request.queryparam.apikey", apiKeyLocation, -1)
	return strings.Replace(tmp, "Verify-API-Key-1", "Verify-API-Key-"+policyName, -1)
}

func AddOAuth2Policy() string {
	return oauth2Policy
}

func AddCORSPolicy() string {
	return corsPolicy
}

func AddSetTargetEndpoint(ref string) string {
	return strings.Replace(setTargetEndpointPolicy, "dynamic.target.url", ref, -1)
}

func replaceTemplateWithPolicy(name string) string {
	re := regexp.MustCompile(`{(.*?)}`)
	return re.ReplaceAllLiteralString(oasPolicyTemplate, name)
}

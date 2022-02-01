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
	<Scope/>
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

var spikeArrestPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<SpikeArrest async="false" continueOnError="false" enabled="true" name="Spike-Arrest-1">    
  <DisplayName>Spike-Arrest-1</DisplayName>
  <Properties/>
  <Rate>1ps</Rate>
  <Identifier/>
  <UseEffectiveCount>true</UseEffectiveCount>
</SpikeArrest>
`

var quotaPolicy1 = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Quota async="false" continueOnError="false" enabled="true" type="calendar" name="Quota-1">
    <DisplayName>Quota-1</DisplayName>
    <Identifier ref="quota.identifier"/>
    <Allow count="1000000000000"/>
    <Interval ref="quota.interval"/>
    <TimeUnit ref="quota.unit"/>
    <Distributed>true</Distributed>
    <StartTime>2019-01-01 00:00:00</StartTime>
</Quota>	
`

var quotaPolicy2 = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Quota async="false" continueOnError="false" enabled="true" type="calendar" name="Quota-1">
    <DisplayName>Quota-1</DisplayName>
    <UseQuotaConfigInAPIProduct>step</UseQuotaConfigInAPIProduct>
    <Distributed>true</Distributed>
    <StartTime>2019-01-01 00:00:00</StartTime>
</Quota>	
`

var graphQLPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<GraphQL async="false" continueOnError="false" enabled="true" name="Validate-name-Schema">
    <Source>request</Source>
    <OperationType>query</OperationType>
    <MaxDepth>4</MaxDepth>
    <MaxCount>4</MaxCount>
    <Action>parse</Action>
    <ResourceURL>graphql://schema.graphql</ResourceURL>
</GraphQL>`

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

func AddSpikeArrestPolicy(policyName string, identifierRef string, rateRef string, rateLiteral string) string {
	policyString := strings.ReplaceAll(spikeArrestPolicy, "Spike-Arrest-1", policyName)
	if rateLiteral != "" {
		rate := "<Rate>" + rateLiteral + "</Rate>"
		policyString = strings.ReplaceAll(policyString, "<Rate>1ps</Rate>", rate)
	} else if rateRef != "" {
		rate := "<Rate ref=\"" + rateRef + "\"/>"
		policyString = strings.ReplaceAll(policyString, "<Rate>1ps</Rate>", rate)
	}
	if identifierRef != "" {
		identifer := "<Identifier ref=\"" + identifierRef + "\"/>"
		policyString = strings.ReplaceAll(policyString, "<Identifier/>", identifer)
	}

	return policyString
}

func AddQuotaPolicy(policyName string, useQuotaConfigStepName string,
	allowRef string, allowLiteral string,
	intervalRef string, intervalLiteral string,
	timeUnitRef string, timeUnitLiteral string) string {
	var policyString string

	if useQuotaConfigStepName != "" {
		policyString = strings.ReplaceAll(quotaPolicy2, "Quota-1", policyName)
		policyString = strings.ReplaceAll(policyString, "step", useQuotaConfigStepName)
	} else {
		policyString = strings.ReplaceAll(quotaPolicy1, "Quota-1", policyName)
		if allowRef != "" {
			allow := "<Allow countRef=\"" + allowRef + "\"/>"
			policyString = strings.ReplaceAll(policyString, "<Allow count=\"1000000000000\"/>", allow)
		} else if allowLiteral != "" {
			allow := "<Allow count=\"" + allowLiteral + "\"/>"
			policyString = strings.ReplaceAll(policyString, "<Allow count=\"1000000000000\"/>", allow)
		}
		if intervalRef != "" {
			interval := "<Interval ref=\"" + intervalRef + "\"/>"
			policyString = strings.ReplaceAll(policyString, "<Interval ref=\"quota.interval\"/>", interval)
		} else if intervalLiteral != "" {
			interval := "<Interval>" + intervalLiteral + "</Interval>"
			policyString = strings.ReplaceAll(policyString, "<Interval ref=\"quota.interval\"/>", interval)
		}
		if timeUnitRef != "" {
			timeUnit := "<TimeUnit ref=\"" + timeUnitRef + "\"/>"
			policyString = strings.ReplaceAll(policyString, "<TimeUnit ref=\"quota.unit\"/>", timeUnit)
		} else if timeUnitLiteral != "" {
			timeUnit := "<TimeUnit>" + timeUnitLiteral + "</TimeUnit>"
			policyString = strings.ReplaceAll(policyString, "<TimeUnit ref=\"quota.unit\"/>", timeUnit)
		}
	}
	return policyString
}

func AddOAuth2Policy(scope string) string {
	if scope != "" {
		scopeTag := "<Scope>" + scope + "</Scope>"
		policyString := strings.ReplaceAll(oauth2Policy, "<Scope/>", scopeTag)
		return policyString
	}
	policyString := strings.ReplaceAll(oauth2Policy, "<Scope/>", "")
	return policyString
}

func AddCORSPolicy() string {
	return corsPolicy
}

func AddSetTargetEndpoint(ref string) string {
	return strings.Replace(setTargetEndpointPolicy, "dynamic.target.url", ref, -1)
}

func AddGraphQLPolicy(name string, action string, schema string) string {
	policyString := strings.ReplaceAll(graphQLPolicy, "schema.graphql", schema)
	policyString = strings.ReplaceAll(policyString, "Validate-name-Schema", "Validate-"+name+"-Schema")
	if action != "" {
		policyString = strings.ReplaceAll(policyString, "parse", action)
	}
	return policyString
}

func replaceTemplateWithPolicy(name string) string {
	re := regexp.MustCompile(`{(.*?)}`)
	return re.ReplaceAllLiteralString(oasPolicyTemplate, name)
}

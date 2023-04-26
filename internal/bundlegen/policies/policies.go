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
	"fmt"
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
    <UseQuotaConfigInAPIProduct stepName="step"/>
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

var setTargetEndpointRefPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<AssignMessage async="false" continueOnError="false" enabled="true" name="Set-Target-1">
    <AssignVariable>
        <Name>target.url</Name>
        <Ref>dynamic.target.url</Ref>
    </AssignVariable>
    <IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
</AssignMessage>`

var setTargetEndpointPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<AssignMessage async="false" continueOnError="false" enabled="true" name="NAME">
	<ReplacePathSuffix/>
	<ReplaceTarget/>
    <IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
	<AssignTo createNew="false" transport="http" type="request"/>
</AssignMessage>`

var setAuthVariablePolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<AssignMessage async="false" continueOnError="false" enabled="true" name="Set-Auth-Var">
    <AssignVariable>
        <Name>auth-var</Name>
        <Value>false</Value>
    </AssignVariable>
    <IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
</AssignMessage>`

var copyAuthHeaderPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<AssignMessage async="false" continueOnError="false" enabled="true" name="Copy-Auth-Var">
	<Set>
		<Headers>
			<Header name="X-Forwarded-Authorization">{request.header.authorization}</Header>
		</Headers>
	</Set>
	<IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
	<AssignTo createNew="false" transport="http" type="request"/>
</AssignMessage>`

var setIntegrationRequestPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<SetIntegrationRequest continueOnError="false" enabled="true" name="set-integration-request">
    <DisplayName>Set Integration Request</DisplayName>
    <ProjectId ref="organization.name"/>
    <IntegrationName>integration_name</IntegrationName>
    <IntegrationRegion ref="system.region.name"/>
    <ApiTrigger>api_trigger/replace_API_1</ApiTrigger>
    <!--
        <Parameters> defines the input parameters to send with the request. Parameters can be int, boolean, String, double, int array, boolean array, String array, double array, and JSON.

        Uncomment the parameters below and modify as needed.
    -->
    <!--
    <Parameters>
        <Parameter name="stringExample" type="string">stringValue</Parameter>
        <Parameter name="doubleExample" type="double">1.0</Parameter>
        <Parameter name="jsonExample" type="json">{}</Parameter>
        <ParameterArray name="intArrayExample" type="integer">
            <Value>1</Value>
            <Value>2</Value>
            <Value>3</Value>
        </ParameterArray>
        <ParameterArray name="booleanArrayExample" type="boolean">
            <Value>true</Value>
            <Value>false</Value>
        </ParameterArray>
    </Parameters>
    -->
</SetIntegrationRequest>`

var verifyJwtPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<VerifyJWT name="SECURITY_POLICY_NAME">
	<DisplayName>SECURITY_POLICY_NAME</DisplayName>
    <Algorithm>RS256</Algorithm>
    <IgnoreUnresolvedVariables>false</IgnoreUnresolvedVariables>
	<Source>request.header.authorization</Source>
    <PublicKey>
        <JWKS uri="JWT_JWKS"/>
    </PublicKey>
    <Issuer>JWT_ISSUER</Issuer>
    <Audience>JWT_AUDIENCE</Audience>
</VerifyJWT>`

var rasiseFaultPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<RaiseFault async="false" continueOnError="false" enabled="true" name="Raise-Fault-Unknown-Request">
    <DisplayName>Raise Fault-Unknown-Request</DisplayName>
    <FaultRules/>
    <Properties/>
    <FaultResponse>
        <Set>
            <Headers/>
            <Payload contentType="application/json" variablePrefix="@" variableSuffix="#">
{
    "error":"invalid_request",
    "error_description": "invalid request"
}
            </Payload>
            <StatusCode>400</StatusCode>
            <ReasonPhrase>Bad Request</ReasonPhrase>
        </Set>
    </FaultResponse>
    <IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
</RaiseFault>`

var extractJwtQueryPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ExtractVariables async="false" continueOnError="false" enabled="true" name="Extract-JWT">
    <DisplayName>Extract-JWT</DisplayName>
	<Source>request</Source>
	<QueryParam name="jwt">
		<Pattern ignoreCase="true">{tokenVar}</Pattern>
 	</QueryParam>
	<IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
</ExtractVariables>`

var extractJwtHeaderPolicy = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ExtractVariables async="false" continueOnError="false" enabled="true" name="Extract-JWT">
    <DisplayName>Extract-JWT</DisplayName>
	<Source>request</Source>
	<Header name="header-name">
		<Pattern ignoreCase="true">value_prefix {tokenVar}</Pattern>
 	</Header>
	<IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
</ExtractVariables>`

var copyAuth = false

func AddSetIntegrationRequestPolicy(integration string, apitrigger string) string {
	policyString := strings.ReplaceAll(setIntegrationRequestPolicy, "integration_name", integration)
	policyString = strings.ReplaceAll(policyString, "replace_API_1", apitrigger)
	return policyString
}

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
	timeUnitRef string, timeUnitLiteral string,
	identifierRef string, identifierLiteral string,
) string {
	var policyString string

	if useQuotaConfigStepName != "" {
		policyString = strings.ReplaceAll(quotaPolicy2, "Quota-1", policyName)
		policyString = strings.ReplaceAll(policyString, "stepName=\"step\"", fmt.Sprintf("stepName=\"%s\"", useQuotaConfigStepName))
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
		if identifierRef == "" && identifierLiteral == "" {
			policyString = strings.ReplaceAll(policyString, "<Identifier ref=\"quota.identifier\"/>", "")
		} else if identifierRef != "" {
			identifier := "<Identifier ref=\"" + identifierRef + "\"/>"
			policyString = strings.ReplaceAll(policyString, "<Identifier ref=\"quota.identifier\"/>", identifier)
		} else if identifierLiteral != "" {
			identifier := "<Identifier>" + identifierLiteral + "</Identifier>"
			policyString = strings.ReplaceAll(policyString, "<Identifier ref=\"quota.identifier\"/>", identifier)
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

func AddSetTargetEndpointRef(ref string) string {
	return strings.Replace(setTargetEndpointRefPolicy, "dynamic.target.url", ref, -1)
}

func AddSetTargetEndpoint(name string, endpoint string, path_transalation string) string {
	pathSuffixTag := `    <AssignVariable>
		<Name>target.copy.pathsuffix</Name>
		<Value>false</Value>
	</AssignVariable>`

	templateSuffixTag := `    <AssignVariable>
		<Name>target.url</Name>
		<Template>{target.url}{proxy.basepath}{proxy.pathsuffix}</Template>
	</AssignVariable>`

	targetUrl := `    <AssignVariable>
		<Name>target.url</Name>
		<Value>VALUE</Value>
	</AssignVariable>`

	policyString := setTargetEndpointPolicy
	if path_transalation == "CONSTANT_ADDRESS" {
		policyString = strings.Replace(policyString, "<ReplacePathSuffix/>", pathSuffixTag, -1)
		policyString = strings.Replace(policyString, "<ReplaceTarget/>", targetUrl, -1)
		policyString = strings.Replace(policyString, "VALUE", endpoint, -1)
	} else {
		policyString = strings.Replace(policyString, "<ReplaceTarget/>", templateSuffixTag, -1)
		policyString = strings.Replace(policyString, "<ReplacePathSuffix/>", "", -1)
	}
	return strings.Replace(policyString, "NAME", name, -1)
}

func AddGraphQLPolicy(name string, action string, schema string) string {
	policyString := strings.ReplaceAll(graphQLPolicy, "schema.graphql", schema)
	policyString = strings.ReplaceAll(policyString, "Validate-name-Schema", "Validate-"+name+"-Schema")
	if action != "" {
		policyString = strings.ReplaceAll(policyString, "parse", action)
	}
	return policyString
}

func AddVerifyJWTPolicy(name string, jwks string, issuer string, audience string, source string) string {
	policyString := strings.ReplaceAll(verifyJwtPolicy, "SECURITY_POLICY_NAME", name)
	policyString = strings.ReplaceAll(policyString, "JWT_JWKS", jwks)
	policyString = strings.ReplaceAll(policyString, "JWT_ISSUER", issuer)
	policyString = strings.ReplaceAll(policyString, "JWT_AUDIENCE", audience)
	if source != "" {
		policyString = strings.ReplaceAll(policyString, "request.header.authorization", source)
	}
	return policyString
}

func AddRaiseFaultPolicy() string {
	return rasiseFaultPolicy
}

func AddExtractJwtQueryPolicy(name string, queryName string) string {
	policyString := strings.ReplaceAll(extractJwtQueryPolicy, "Extract-JWT", name)
	policyString = strings.ReplaceAll(policyString, "jwt", queryName)
	return policyString
}

func AddExtractJwtHeaderPolicy(name string, headerName string, prefixName string) string {
	policyString := strings.ReplaceAll(extractJwtHeaderPolicy, "Extract-JWT", name)
	policyString = strings.ReplaceAll(policyString, "header-name", headerName)
	policyString = strings.ReplaceAll(policyString, "value_prefix", prefixName)
	return policyString
}

// TODO: Unused at the moment
func AddSetAuthVarPolicy(auth bool) string {
	policyString := strings.ReplaceAll(setAuthVariablePolicy, "<Value>false</Value>", fmt.Sprintf("<Value>%t</Value>", auth))
	return policyString
}

func AddCopyAuthHeaderPolicy() string {
	return copyAuthHeaderPolicy
}

func EnableCopyAuthPolicy() {
	copyAuth = true
}

func IsCopyAuthEnabled() bool {
	return copyAuth
}

func replaceTemplateWithPolicy(name string) string {
	re := regexp.MustCompile(`{(.*?)}`)
	return re.ReplaceAllLiteralString(oasPolicyTemplate, name)
}

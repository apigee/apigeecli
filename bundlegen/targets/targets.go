// Copyright 2021 Google LLC
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

package targets

import (
	"encoding/xml"

	proxytypes "github.com/apigee/apigeecli/bundlegen/common"
)

type targetEndpointDef struct {
	XMLName              xml.Name                `xml:"TargetEndpoint"`
	Name                 string                  `xml:"name,attr"`
	FaultRules           string                  `xml:"FaultRules,omitempty"`
	PreFlow              proxytypes.PreFlowDef   `xml:"PreFlow,omitempty"`
	PostFlow             proxytypes.PostFlowDef  `xml:"PostFlow,omitempty"`
	Flows                proxytypes.FlowsDef     `xml:"Flows,omitempty"`
	HTTPTargetConnection httpTargetConnectionDef `xml:"HTTPTargetConnection,omitempty"`
}

type property struct {
	XMLName xml.Name `xml:"Property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

type properties struct {
	XMLName  xml.Name   `xml:"Properties"`
	Property []property `xml:"Property"`
}

type httpTargetConnectionDef struct {
	Authentication *authenticationDef `xml:"Authentication"`
	URL            string             `xml:"URL"`
	Properties     properties         `xml:"Properties"`
}

type authenticationDef struct {
	XMLName           xml.Name              `xml:"Authentication"`
	GoogleAccessToken *googleAccessTokenDef `xml:"GoogleAccessToken,omitempty"`
	GoogleIDToken     *googleIdTokenDef     `xml:"GoogleIDToken,omitempty"`
}

type googleAccessTokenDef struct {
	XMLName xml.Name   `xml:"GoogleAccessToken"`
	Scopes  []scopeDef `xml:"Scopes,omitempty"`
}

type googleIdTokenDef struct {
	XMLName  xml.Name     `xml:"GoogleIDToken"`
	Audience *audienceDef `xml:"Audience"`
}

type audienceDef struct {
	XMLName      xml.Name `xml:"Audience"`
	Ref          *string  `xml:"ref,attr"`
	UseTargetUrl *string  `xml:"useTargetUrl,attr"`
}

type scopeDef struct {
	Scope string `xml:"Scope"`
}

var integrationEndpoint = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<IntegrationEndpoint name="default">
    <AsyncExecution>false</AsyncExecution>
</IntegrationEndpoint>
`

var TargetEndpoints []targetEndpointDef

func AddStepToPreFlowRequest(name string, targetEndpointName string) {
	for _, targetEndpoint := range TargetEndpoints {
		if targetEndpoint.Name == targetEndpointName {
			step := proxytypes.StepDef{}
			step.Name = name
			targetEndpoint.PreFlow.Request.Step = append(targetEndpoint.PreFlow.Request.Step, &step)
		}
	}
}

func GetTargetEndpoint(targetEndpoint targetEndpointDef) (string, error) {
	targetBody, err := xml.MarshalIndent(targetEndpoint, "", " ")
	if err != nil {
		return "", nil
	}
	return string(targetBody), nil
}

func NewTargetEndpoint(name string, endpoint string, oasGoogleAcessTokenScopeLiteral string, oasGoogleIdTokenAudLiteral string, oasGoogleIdTokenAudRef string) {
	targetEndpoint := targetEndpointDef{}
	targetEndpoint.Name = name
	targetEndpoint.PreFlow.Name = "PreFlow"
	targetEndpoint.PostFlow.Name = "PostFlow"
	targetEndpoint.HTTPTargetConnection.URL = endpoint

	if oasGoogleAcessTokenScopeLiteral != "" {
		targetEndpoint.HTTPTargetConnection.Authentication = new(authenticationDef)
		targetEndpoint.HTTPTargetConnection.Authentication.GoogleAccessToken = new(googleAccessTokenDef)
		targetEndpoint.HTTPTargetConnection.Authentication.GoogleAccessToken.Scopes = make([]scopeDef, 1)
		targetEndpoint.HTTPTargetConnection.Authentication.GoogleAccessToken.Scopes[0].Scope = oasGoogleAcessTokenScopeLiteral

	} else if oasGoogleIdTokenAudLiteral != "" || oasGoogleIdTokenAudRef != "" {
		targetEndpoint.HTTPTargetConnection.Authentication = new(authenticationDef)
		targetEndpoint.HTTPTargetConnection.Authentication.GoogleIDToken = new(googleIdTokenDef)

		targetEndpoint.HTTPTargetConnection.Authentication.GoogleIDToken.Audience = new(audienceDef)

		if oasGoogleIdTokenAudLiteral != "" {
			targetEndpoint.HTTPTargetConnection.Authentication.GoogleIDToken.Audience.XMLName.Local = oasGoogleIdTokenAudLiteral
		} else if oasGoogleIdTokenAudRef != "" {
			targetEndpoint.HTTPTargetConnection.Authentication.GoogleIDToken.Audience.Ref = new(string)
			targetEndpoint.HTTPTargetConnection.Authentication.GoogleIDToken.Audience.Ref = setString(oasGoogleIdTokenAudRef)
			targetEndpoint.HTTPTargetConnection.Authentication.GoogleIDToken.Audience.UseTargetUrl = new(string)
			*targetEndpoint.HTTPTargetConnection.Authentication.GoogleIDToken.Audience.UseTargetUrl = "true"
		}
	} else {
		targetEndpoint.HTTPTargetConnection.Authentication = nil
	}
	TargetEndpoints = append(TargetEndpoints, targetEndpoint)
}

func IsExists(endpointName string) bool {
	for index := range TargetEndpoints {
		if TargetEndpoints[index].Name == endpointName {
			return true
		}
	}
	return false
}

func AddTargetEndpointProperty(endpointName string, propertyName string, propertyValue string) {

	property := property{}
	property.Name = propertyName
	property.Value = propertyValue

	for index := range TargetEndpoints {
		if TargetEndpoints[index].Name == endpointName {
			TargetEndpoints[index].HTTPTargetConnection.Properties.Property = append(TargetEndpoints[index].HTTPTargetConnection.Properties.Property, property)
			return
		}
	}
}

func GetIntegrationEndpoint() string {
	return integrationEndpoint
}

func setString(s string) *string {
	return &s
}

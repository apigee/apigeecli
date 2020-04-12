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

package proxies

import (
	"encoding/xml"
)

type proxyEndpointDef struct {
	XMLName             xml.Name               `xml:"ProxyEndpoint"`
	Name                string                 `xml:"name,attr"`
	Description         string                 `xml:"Description,omitempty"`
	FaultRules          string                 `xml:"FaultRules,omitempty"`
	PreFlow             preFlowDef             `xml:"PreFlow,omitempty"`
	PostFlow            postFlowDef            `xml:"PostFlow,omitempty"`
	Flows               flowsDef               `xml:"Flows,omitempty"`
	HTTPProxyConnection httpProxyConnectionDef `xml:"HTTPProxyConnection,omitempty"`
	RouteRule           routeRuleDef           `xml:"RouteRule,omitempty"`
}

type preFlowDef struct {
	XMLName  xml.Name        `xml:"PreFlow"`
	Name     string          `xml:"name,attr"`
	Request  requestFlowDef  `xml:"Request"`
	Response responseFlowDef `xml:"Response"`
}

type postFlowDef struct {
	XMLName  xml.Name        `xml:"PostFlow"`
	Name     string          `xml:"name,attr"`
	Request  requestFlowDef  `xml:"Request"`
	Response responseFlowDef `xml:"Response"`
}

type requestFlowDef struct {
	Step []*stepDef `xml:"Step"`
}

type responseFlowDef struct {
	Step []*stepDef `xml:"Step"`
}

type stepDef struct {
	Name string `xml:"Name"`
}

type routeRuleDef struct {
	XMLName        xml.Name `xml:"RouteRule"`
	Name           string   `xml:"name,attr"`
	TargetEndpoint string   `xml:"TargetEndpoint"`
}

type flowsDef struct {
	XMLName xml.Name  `xml:"Flows"`
	Flow    []flowDef `xml:"Flow"`
}

type flowDef struct {
	XMLName     xml.Name        `xml:"Flow"`
	Name        string          `xml:"name,attr"`
	Description string          `xml:"Description,omitempty"`
	Request     requestFlowDef  `xml:"Request"`
	Response    responseFlowDef `xml:"Response"`
	Condition   conditionDef    `xml:"Condition"`
}

type conditionDef struct {
	ConditionData string `xml:",innerxml"`
}

type httpProxyConnectionDef struct {
	XMLName     xml.Name `xml:"HTTPProxyConnection"`
	BasePath    string   `xml:"BasePath"`
	Properties  string   `xml:"Properties"`
	VirtualHost []string `xml:"VirtualHost"`
}

var proxyEndpoint proxyEndpointDef

func GetProxyEndpoint() (string, error) {
	proxyBody, err := xml.MarshalIndent(proxyEndpoint, "", " ")
	if err != nil {
		return "", nil
	}
	return string(proxyBody), nil
}

func NewProxyEndpoint(basePath string) {
	proxyEndpoint.Name = "default"
	proxyEndpoint.PreFlow.Name = "PreFlow"
	proxyEndpoint.PostFlow.Name = "PostFlow"
	proxyEndpoint.HTTPProxyConnection.BasePath = basePath
	proxyEndpoint.RouteRule.Name = "default"
	proxyEndpoint.RouteRule.TargetEndpoint = "default"
}

func AddFlow(operationId string, keyPath string, method string, description string) {
	flow := flowDef{}
	flow.Name = operationId
	flow.Description = description
	flow.Condition.ConditionData = "(proxy.pathsuffix MatchesPath \"" + keyPath + "\") and (request.verb = \"" + method + "\")"
	proxyEndpoint.Flows.Flow = append(proxyEndpoint.Flows.Flow, flow)
}

func AddStepToPreFlowRequest(name string) {
	step := stepDef{}
	step.Name = name
	proxyEndpoint.PreFlow.Request.Step = append(proxyEndpoint.PreFlow.Request.Step, &step)
}

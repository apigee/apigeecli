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

package proxies

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	proxytypes "internal/bundlegen/common"
)

type proxyEndpointDef struct {
	XMLName             xml.Name               `xml:"ProxyEndpoint"`
	Name                string                 `xml:"name,attr"`
	Description         string                 `xml:"Description,omitempty"`
	FaultRules          string                 `xml:"FaultRules,omitempty"`
	PreFlow             proxytypes.PreFlowDef  `xml:"PreFlow,omitempty"`
	PostFlow            proxytypes.PostFlowDef `xml:"PostFlow,omitempty"`
	Flows               proxytypes.FlowsDef    `xml:"Flows,omitempty"`
	HTTPProxyConnection httpProxyConnectionDef `xml:"HTTPProxyConnection,omitempty"`
	RouteRule           []routeRuleDef         `xml:"RouteRule,omitempty"`
}

type routeRuleDef struct {
	XMLName             xml.Name `xml:"RouteRule"`
	Name                string   `xml:"name,attr"`
	Condition           *string  `xml:"Condition"`
	Url                 *string  `xml:"URL"`
	TargetEndpoint      *string  `xml:"TargetEndpoint"`
	IntegrationEndpoint *string  `xml:"IntegrationEndpoint"`
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
	proxyBody = bytes.ReplaceAll(proxyBody, []byte("&#34;"), []byte("\""))
	return string(proxyBody), nil
}

func NewProxyEndpoint(basePath string, targetEndpoint bool) {
	routeRule := routeRuleDef{}
	proxyEndpoint.Name = "default"
	proxyEndpoint.PreFlow.Name = "PreFlow"
	proxyEndpoint.PostFlow.Name = "PostFlow"
	proxyEndpoint.HTTPProxyConnection.BasePath = basePath
	routeRule.Name = "default"
	if targetEndpoint {
		routeRule.TargetEndpoint = new(string)
		*routeRule.TargetEndpoint = "default"
	} else {
		routeRule.IntegrationEndpoint = new(string)
		*routeRule.IntegrationEndpoint = "default"
	}
	proxyEndpoint.RouteRule = append(proxyEndpoint.RouteRule, routeRule)
}

func AddFlow(operationId string, keyPath string, method string, description string) {
	flow := proxytypes.FlowDef{}
	flow.Name = operationId
	if description != "" {
		flow.Description = description
	}
	if keyPath != "" && method != "" {
		flow.Condition.ConditionData = "(proxy.pathsuffix MatchesPath \"" + keyPath + "\") and (request.verb = \"" + strings.ToUpper(method) + "\")"
	}
	proxyEndpoint.Flows.Flow = append(proxyEndpoint.Flows.Flow, flow)
}

func FlowExists(name string) bool {
	for _, flow := range proxyEndpoint.Flows.Flow {
		if flow.Name == name {
			return true
		}
	}
	return false
}

func AddStepToPreFlowRequest(name string) {
	step := proxytypes.StepDef{}
	step.Name = name
	proxyEndpoint.PreFlow.Request.Step = append(proxyEndpoint.PreFlow.Request.Step, &step)
}

func AddStepToFlowRequest(name string, flowName string) error {
	for flowKey, flow := range proxyEndpoint.Flows.Flow {
		if flow.Name == flowName {
			step := proxytypes.StepDef{}
			step.Name = name
			proxyEndpoint.Flows.Flow[flowKey].Request.Step = append(proxyEndpoint.Flows.Flow[flowKey].Request.Step, &step)
			return nil
		}
	}
	return fmt.Errorf("flow name not found")
}

func AddRoute(name string, endpoint string, condition string) {
	routeRule := routeRuleDef{}
	routeRule.Name = name

	routeRule.TargetEndpoint = new(string)
	*routeRule.TargetEndpoint = endpoint

	routeRule.Condition = new(string)
	*routeRule.Condition = condition
	proxyEndpoint.RouteRule = append(proxyEndpoint.RouteRule, routeRule)
}

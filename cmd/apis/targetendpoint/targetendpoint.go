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

package targetendpoint

import (
	"encoding/xml"
)

type preFlowDef struct {
	XMLName  xml.Name        `xml:"PreFlow"`
	Name     string          `xml:"name,attr"`
	Request  requestFlowDef  `xml:"Request,omitempty"`
	Response responseFlowDef `xml:"Response,omitempty"`
}

type postFlowDef struct {
	XMLName  xml.Name        `xml:"PostFlow"`
	Name     string          `xml:"name,attr"`
	Request  requestFlowDef  `xml:"Request,omitempty"`
	Response responseFlowDef `xml:"Response,omitempty"`
}

type httpTargetConnectionDef struct {
	URL string `xml:"URL"`
}

type requestFlowDef struct {
	Step []*stepDef `xml:"Step,omitempty"`
}

type responseFlowDef struct {
	Step []*stepDef `xml:"Step,omitempty"`
}

type stepDef struct {
	Name string `xml:"Name"`
}

type targetEndpointDef struct {
	XMLName              xml.Name                `xml:"TargetEndpoint"`
	Name                 string                  `xml:"name,attr"`
	PreFlow              preFlowDef              `xml:"PreFlow,omitempty"`
	PostFlow             postFlowDef             `xml:"PostFlow,omitempty"`
	HTTPTargetConnection httpTargetConnectionDef `xml:"HTTPTargetConnection,omitempty"`
}

var targetEndpoint targetEndpointDef

func GetTargetEndpoint() (string, error) {
	targetBody, err := xml.MarshalIndent(targetEndpoint, "", " ")
	if err != nil {
		return "", nil
	}
	return string(targetBody), nil
}

func NewTargetEndpoint(endpoint string) {
	targetEndpoint.Name = "default"
	targetEndpoint.PreFlow.Name = "PreFlow"
	targetEndpoint.PostFlow.Name = "PostFlow"
	targetEndpoint.HTTPTargetConnection.URL = endpoint
}

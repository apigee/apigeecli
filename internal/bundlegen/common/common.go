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

package common

import (
	"encoding/xml"
)

type PreFlowDef struct {
	XMLName  xml.Name        `xml:"PreFlow"`
	Name     string          `xml:"name,attr"`
	Request  RequestFlowDef  `xml:"Request"`
	Response ResponseFlowDef `xml:"Response"`
}

type PostFlowDef struct {
	XMLName  xml.Name        `xml:"PostFlow"`
	Name     string          `xml:"name,attr"`
	Request  RequestFlowDef  `xml:"Request"`
	Response ResponseFlowDef `xml:"Response"`
}

type RequestFlowDef struct {
	Step []*StepDef `xml:"Step"`
}

type ResponseFlowDef struct {
	Step []*StepDef `xml:"Step"`
}

type StepDef struct {
	Name string `xml:"Name"`
}

type FlowsDef struct {
	XMLName xml.Name  `xml:"Flows"`
	Flow    []FlowDef `xml:"Flow"`
}

type FlowDef struct {
	XMLName     xml.Name        `xml:"Flow"`
	Name        string          `xml:"name,attr"`
	Description string          `xml:"Description,omitempty"`
	Request     RequestFlowDef  `xml:"Request"`
	Response    ResponseFlowDef `xml:"Response"`
	Condition   ConditionDef    `xml:"Condition"`
}

type ConditionDef struct {
	ConditionData string `xml:",innerxml"`
}

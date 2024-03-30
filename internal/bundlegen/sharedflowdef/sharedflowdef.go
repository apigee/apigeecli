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

package sharedflowdef

import "encoding/xml"

type sharedFlowDef struct {
	XMLName        xml.Name       `xml:"SharedFlowBundle"`
	Name           string         `xml:"name,attr"`
	Revision       string         `xml:"revision,attr"`
	CreatedAt      string         `xml:"CreatedAt,omitempty"`
	Description    string         `xml:"Description,omitempty"`
	DisplayName    string         `xml:"DisplayName,omitempty"`
	LastModifiedAt string         `xml:"LastModifiedAt,omitempty"`
	Policies       policiesDef    `xml:"Policies,omitempty"`
	SharedFlows    sharedflowsDef `xml:"SharedFlows,omitempty"`
	Resources      resourcesDef   `xml:"Resources,omitempty"`
	SubType        string         `xml:"subType,omitempty"`
}

type resourcesDef struct {
	Resource []string `xml:"Resource,omitempty"`
}

type policiesDef struct {
	Policy []string `xml:"Policy,omitempty"`
}

type sharedflowsDef struct {
	SharedFlow []string `xml:"SharedFlow,omitempty"`
}

func SetDescriptionWithMarshal(contents []byte, description string) (resp []byte, err error) {
	s := sharedFlowDef{}
	err = xml.Unmarshal(contents, &s)
	if err != nil {
		return nil, err
	}
	if s.Description != "" {
		s.Description = s.Description + " " + description
	} else {
		s.Description = description
	}
	resp, err = xml.MarshalIndent(s, "", " ")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

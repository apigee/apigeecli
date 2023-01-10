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

package apiproxydef

import (
	"encoding/xml"
	"strconv"
	"time"
)

type policiesDef struct {
	Policy []string `xml:"Policy,omitempty"`
}

type proxyEndpointsDef struct {
	XMLName       xml.Name `xml:"ProxyEndpoints"`
	ProxyEndpoint []string `xml:"ProxyEndpoint,omitempty"`
}

type targetEndpointsDef struct {
	XMLName        xml.Name `xml:"TargetEndpoints"`
	TargetEndpoint []string `xml:"TargetEndpoint,omitempty"`
}

type integrationEndpointsDef struct {
	XMLName             xml.Name `xml:"IntegrationEndpoints"`
	IntegrationEndpoint []string `xml:"IntegrationEndpoint,omitempty"`
}

type configurationVersionDef struct {
	XMLName      xml.Name `xml:"ConfigurationVersion,omitempty"`
	MajorVersion string   `xml:"majorVersion,attr"`
	MinorVersion string   `xml:"minorVersion,attr"`
}

type resourcesDef struct {
	Resource []string `xml:"Resource,omitempty"`
}

type apiProxyDef struct {
	XMLName              xml.Name                `xml:"APIProxy"`
	Name                 string                  `xml:"name,attr"`
	Revision             string                  `xml:"revision,attr"`
	BasePaths            string                  `xml:"Basepaths,omitempty"`
	ConfigurationVersion configurationVersionDef `xml:"ConfigurationVersion,omitempty"`
	CreatedAt            string                  `xml:"CreatedAt,omitempty"`
	Description          string                  `xml:"Description,omitempty"`
	DisplayName          string                  `xml:"DisplayName,omitempty"`
	LastModifiedAt       string                  `xml:"LastModifiedAt,omitempty"`
	Policies             policiesDef             `xml:"Policies,omitempty"`
	ProxyEndpoints       proxyEndpointsDef       `xml:"ProxyEndpoints,omitempty"`
	IntegrationEndpoints integrationEndpointsDef `xml:"IntegrationEndpoints,omitempty"`
	Resources            resourcesDef            `xml:"Resources,omitempty"`
	Spec                 string                  `xml:"Spec,omitempty"`
	TargetServers        string                  `xml:"TargetServers,omitempty"`
	TargetEndpoints      targetEndpointsDef      `xml:"TargetEndpoints,omitempty"`
	Validate             string                  `xml:"validate,omitempty"`
}

var apiProxy apiProxyDef

func SetDisplayName(name string) {
	apiProxy.DisplayName = name
	apiProxy.Name = name
}

func AddProxyEndpoint(name string) {
	apiProxy.ProxyEndpoints.ProxyEndpoint = append(apiProxy.ProxyEndpoints.ProxyEndpoint, name)
}

func AddTargetEndpoint(name string) {
	apiProxy.TargetEndpoints.TargetEndpoint = append(apiProxy.TargetEndpoints.TargetEndpoint, name)
}

func AddIntegrationEndpoint(name string) {
	apiProxy.IntegrationEndpoints.IntegrationEndpoint = append(apiProxy.IntegrationEndpoints.IntegrationEndpoint, name)
}

func SetCreatedAt() {
	apiProxy.CreatedAt = strconv.FormatInt((time.Now().UTC().UnixNano())/1000000, 10)
}

func SetLastModifiedAt() {
	apiProxy.LastModifiedAt = strconv.FormatInt((time.Now().UTC().UnixNano())/1000000, 10)
}

func AddPolicy(name string) {
	for index := range apiProxy.Policies.Policy {
		if apiProxy.Policies.Policy[index] == name {
			return
		}
	}
	apiProxy.Policies.Policy = append(apiProxy.Policies.Policy, name)
}

func SetBasePath(basePath string) {
	apiProxy.BasePaths = basePath
}

func SetRevision(revision string) {
	apiProxy.Revision = revision
}

func SetDescription(description string) {
	apiProxy.Description = description
}

func GetAPIProxy() (string, error) {
	proxyBody, err := xml.MarshalIndent(apiProxy, "", " ")
	if err != nil {
		return "", err
	}
	return string(proxyBody), nil
}

func SetConfigurationVersion() {
	apiProxy.ConfigurationVersion.MajorVersion = "4"
	apiProxy.ConfigurationVersion.MinorVersion = "0"
}

func AddResource(name string, resType string) {
	apiProxy.Resources.Resource = append(apiProxy.Resources.Resource, resType+"://"+name)
}

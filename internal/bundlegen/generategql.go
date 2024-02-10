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

package bundlegen

import (
	"fmt"
	"net/url"

	apiproxy "internal/bundlegen/apiproxydef"
	"internal/bundlegen/proxies"
	targets "internal/bundlegen/targets"
)

func GenerateAPIProxyDefFromGQL(name string,
	gqlDocName string,
	basePath string,
	apiKeyLocation string,
	skipPolicy bool,
	addCORS bool,
	targetUrlRef string,
	targetUrl string,
) (err error) {
	apiproxy.SetDisplayName(name)
	apiproxy.SetCreatedAt()
	apiproxy.SetLastModifiedAt()
	apiproxy.SetConfigurationVersion()
	apiproxy.AddTargetEndpoint(DEFAULT)
	apiproxy.AddProxyEndpoint(DEFAULT)

	apiproxy.SetDescription("Generated API Proxy from " + gqlDocName)

	if !skipPolicy {
		apiproxy.AddResource(gqlDocName, "graphql")
		apiproxy.AddPolicy("Validate-" + name + "-Schema")
	}

	proxies.NewProxyEndpoint(basePath, true)

	if addCORS {
		proxies.AddStepToPreFlowRequest("Add-CORS")
		apiproxy.AddPolicy("Add-CORS")
	}

	// if target is not set, add a default/fake endpoint
	if targetUrl == "" {
		targets.NewTargetEndpoint(DEFAULT, "https://api.example.com", "", "", "", "", "")
	} else { // an explicit target url is set
		if _, err = url.Parse(targetUrl); err != nil {
			return fmt.Errorf("invalid target url: %v", err)
		}
		targets.NewTargetEndpoint(DEFAULT, targetUrl, "", "", "", "", "")
	}

	// set a dynamic target url
	if targetUrlRef != "" {
		targets.AddStepToPreFlowRequest("Set-Target-1", DEFAULT)
		apiproxy.AddPolicy("Set-Target-1")
		generateSetTarget = true
	}

	if !skipPolicy {
		proxies.AddStepToPreFlowRequest("Validate-" + name + "-Schema")
	}

	if apiKeyLocation != "" {
		proxies.AddStepToPreFlowRequest("Verify-API-Key-" + name)
	}

	return err
}

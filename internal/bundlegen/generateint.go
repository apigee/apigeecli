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
	apiproxy "internal/bundlegen/apiproxydef"
	"internal/bundlegen/proxies"
)

func GenerateIntegrationAPIProxy(name string,
	apitrigger string,
) (err error) {
	apiproxy.SetDisplayName(name)
	apiproxy.SetCreatedAt()
	apiproxy.SetLastModifiedAt()
	apiproxy.SetConfigurationVersion()
	apiproxy.AddProxyEndpoint("default")
	apiproxy.AddIntegrationEndpoint("default")
	apiproxy.SetBasePath("/" + apitrigger)

	proxies.NewProxyEndpoint("/"+apitrigger, false)

	proxies.AddStepToPreFlowRequest("set-integration-request")
	apiproxy.AddPolicy("set-integration-request")

	return nil
}

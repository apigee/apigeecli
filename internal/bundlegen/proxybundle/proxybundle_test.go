// Copyright 2024 Google LLC
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

package proxybundle

import (
	"os"
	"path"
	"regexp"
	"testing"

	"internal/bundlegen"
	"internal/client/clienttest"
)

var specName = "petstore-v3.1.yaml"

var cliPath = os.Getenv("APIGEECLI_PATH")

func TestGenerateAPIProxyBundleFromOAS(t *testing.T) {
	var contents []byte
	var err error

	if err = clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if contents, err = bundlegen.LoadDocument(path.Join(cliPath, "test"), "", specName, true); err != nil {
		t.Fatalf("%v", err)
	}
	skipPolicy := false
	name := "test"
	addCORS := true
	basePath := ""
	oasGoogleAcessTokenScopeLiteral := ""
	targetURL := "http://api.example.com"
	oasGoogleIDTokenAudLiteral := ""
	oasGoogleIDTokenAudRef := ""
	targetURLRef := ""
	targetServerName := ""
	version := bundlegen.GetModelVersion()
	if version != "" {
		re := regexp.MustCompile(`3\.1\.[0-9]`)
		if re.MatchString(version) {
			skipPolicy = true
		}
	}
	targetOptions := bundlegen.TargetOptions{
		IntegrationBackend: bundlegen.IntegrationBackendOptions{
			IntegrationName: "",
			TriggerName:     "",
		},
		HttpBackend: bundlegen.HttpBackendOptions{
			OasGoogleAcessTokenScopeLiteral: oasGoogleAcessTokenScopeLiteral,
			OasGoogleIDTokenAudLiteral:      oasGoogleIDTokenAudLiteral,
			OasGoogleIDTokenAudRef:          oasGoogleIDTokenAudRef,
			OasTargetURLRef:                 targetURLRef,
			TargetURL:                       targetURL,
			TargetServerName:                targetServerName,
		},
	}
	if err := bundlegen.GenerateAPIProxyDefFromOASv2(name,
		basePath,
		specName,
		skipPolicy,
		addCORS,
		targetOptions); err != nil {
		t.Fatalf("%v", err)
	}
	// Create the API proxy bundle
	if err := GenerateAPIProxyBundleFromOAS(name,
		string(contents),
		specName,
		skipPolicy,
		addCORS,
		targetOptions); err != nil {
		t.Fatalf("%v", err)
	}
}

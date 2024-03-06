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

package bundlegen

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"testing"

	"internal/client/clienttest"
)

var cliPath = os.Getenv("APIGEECLI_PATH")

var specNames = []string{
	"petstore-v3.1.yaml",
	"md-ext1.yaml",
	"windfarm3-quota.yaml",
	"open_weather_api.yaml",
}

func TestSpecs(t *testing.T) {
	for _, specName := range specNames {
		fmt.Println("Testing " + specName + " ...")
		testLoadDocument(specName, t)
		TestGenerateAPIProxyDefFromOASv2(specName, t)
	}
}

func testLoadDocument(specName string, t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := LoadDocument(path.Join(cliPath, "test"), "", specName, true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGenerateAPIProxyDefFromOASv2(specName string, t *testing.T) {
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
	version := GetModelVersion()
	if version != "" {
		re := regexp.MustCompile(`3\.1\.[0-9]`)
		if re.MatchString(version) {
			skipPolicy = true
		}
	}

	targetOptions := TargetOptions{
		IntegrationBackend: IntegrationBackendOptions{},
		HttpBackend: HttpBackendOptions{
			OasGoogleAcessTokenScopeLiteral: oasGoogleAcessTokenScopeLiteral,
			OasGoogleIDTokenAudLiteral:      oasGoogleIDTokenAudLiteral,
			OasGoogleIDTokenAudRef:          oasGoogleIDTokenAudRef,
			OasTargetURLRef:                 targetURLRef,
			TargetURL:                       targetURL,
			TargetServerName:                targetServerName,
		},
	}

	if err := GenerateAPIProxyDefFromOASv2(name,
		basePath,
		specName,
		skipPolicy,
		addCORS,
		targetOptions); err != nil {
		t.Fatalf("%v", err)
	}
}

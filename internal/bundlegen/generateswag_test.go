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
	"internal/client/clienttest"
	"path"
	"testing"
)

var swagSpecNames = []string{
	"dynamicroute_swagger.json",
	"endpoints1.yaml",
	"endpoints2.yaml",
}

var oasDocName string

func TestSwagSpecs(t *testing.T) {
	for _, specName := range swagSpecNames {
		fmt.Println("Testing " + specName + " ...")
		testLoadSwaggerFromFile(specName, t)
		testGenerateAPIProxyFromSwagger(t)
	}
}

func testLoadSwaggerFromFile(specName string, t *testing.T) {
	var err error
	if err = clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if oasDocName, _, err = LoadSwaggerFromFile(path.Join(cliPath, "test", specName)); err != nil {
		t.Fatalf("%v", err)
	}
}

func testGenerateAPIProxyFromSwagger(t *testing.T) {
	name := "test"
	basePath := ""
	addCORS := true
	if _, err := GenerateAPIProxyFromSwagger(name, "desc", oasDocName, basePath, addCORS); err != nil {
		t.Fatalf("%v", err)
	}
}

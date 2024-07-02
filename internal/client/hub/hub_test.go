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

package hub

import (
	"os"
	"path"
	"testing"

	"internal/client/clienttest"
	"internal/cmd/utils"
)

var cliPath = os.Getenv("APIGEECLI_PATH")

func TestCreateApi(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	apiID := "test-api"
	var contents []byte

	if contents, err = utils.ReadFile(path.Join(cliPath, "test", "api.json")); err != nil {
		t.Errorf("failed to read api.json: %v", err)
	}

	if _, err = CreateApi(apiID, contents); err != nil {
		t.Errorf("failed to create api: %v", err)
	}
}

func TestGetApi(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	apiID := "test-api"
	if _, err = GetApi(apiID); err != nil {
		t.Errorf("failed to get api: %v", err)
	}
}

func TestListApi(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	if _, err = ListApi("", -1, ""); err != nil {
		t.Errorf("failed to list api: %v", err)
	}
}

func TestCreateApiVersion(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	var contents []byte

	if contents, err = utils.ReadFile(path.Join(cliPath, "test", "api-ver.json")); err != nil {
		t.Errorf("failed to read api.json: %v", err)
	}

	apiID := "test-api"
	versionID := "test-version"

	if _, err = CreateApiVersion(apiID, versionID, contents); err != nil {
		t.Errorf("failed to create api version: %v", err)
	}
}

func TestGetApiVersion(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	apiID := "test-api"
	versionID := "test-version"

	if _, err = GetApiVersion(apiID, versionID); err != nil {
		t.Errorf("failed to get api version: %v", err)
	}
}

func TestListApiVersion(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	apiID := "test-api"

	if _, err = ListApiVersions(apiID, "", -1, ""); err != nil {
		t.Errorf("failed to list api version: %v", err)
	}
}

func TestCreateApiVersionsSpec(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	var contents []byte

	if contents, err = utils.ReadFile(path.Join(cliPath, "test", "petstore-v3.1.json")); err != nil {
		t.Errorf("failed to read api.json: %v", err)
	}

	apiID := "test-api"
	versionID := "test-version"
	specID := "test-spec"
	displayName := "test-spec"
	mimeType := "json"
	sourceURI := ""
	documentation := ""

	if _, err = CreateApiVersionsSpec(apiID, versionID, specID, displayName, contents, mimeType, sourceURI, documentation); err != nil {
		t.Errorf("failed to create api version spec: %v", err)
	}
}

func TestGetApiVersionSpec(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	apiID := "test-api"
	versionID := "test-version"
	specID := "test-spec"

	if _, err = GetApiVersionSpec(apiID, versionID, specID); err != nil {
		t.Errorf("failed to get api version spec: %v", err)
	}
}

func TestGetApiVersionsSpecContents(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	apiID := "test-api"
	versionID := "test-version"
	specID := "test-spec"

	if _, err = GetApiVersionsSpecContents(apiID, versionID, specID); err != nil {
		t.Errorf("failed to get api version spec contents: %v", err)
	}
}

func TestLintApiVersionSpecs(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	apiID := "test-api"
	versionID := "test-version"
	specID := "test-spec"

	if _, err = LintApiVersionSpec(apiID, versionID, specID); err != nil {
		t.Errorf("failed to lint api version spec: %v", err)
	}
}

func TestListApiVersionSpecs(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	apiID := "test-api"
	versionID := "test-version"

	if _, err = ListApiVersionSpecs(apiID, versionID, "", -1, ""); err != nil {
		t.Errorf("failed to list api version specs: %v", err)
	}
}

func TestDeleteApiVersionSpec(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	apiID := "test-api"
	versionID := "test-version"
	specID := "test-spec"

	if _, err = DeleteApiVersionSpec(apiID, versionID, specID); err != nil {
		t.Errorf("failed to delete api version spec: %v", err)
	}
}

func TestDeleteApiVersion(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	apiID := "test-api"
	versionID := "test-version"

	if _, err = DeleteApiVersion(apiID, versionID); err != nil {
		t.Errorf("failed to delete api version: %v", err)
	}
}

func TestDeleteApi(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	apiID := "test-api"

	if _, err = DeleteApi(apiID); err != nil {
		t.Errorf("failed to delete api: %v", err)
	}
}

func TestCreateAttribute(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	var aValues []byte
	attributeID := "test-attribute"
	displayName := "test attribute"
	description := "test description"
	scope := "API"
	dataType := "ENUM"
	cardinality := 1

	if _, err := CreateAttribute(attributeID, displayName, description,
		scope, dataType, aValues, cardinality); err != nil {
		t.Errorf("failed to create attribute %v", err)
	}
}

func TestGetAttribute(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	attributeID := "test-attribute"

	if _, err := GetAttribute(attributeID); err != nil {
		t.Errorf("failed to get attribute %v", err)
	}
}

func TestListAttribute(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	if _, err := ListAttributes("", -1, ""); err != nil {
		t.Errorf("failed to list attributes %v", err)
	}
}

func TestDeleteAttribute(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	attributeID := "test-attribute"

	if _, err := DeleteAttribute(attributeID); err != nil {
		t.Errorf("failed to get attribute %v", err)
	}
}

func TestCreateDependency(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	dependencyID := "test-dependency"
	description := "test description"
	consumerDisplayName := "test consumer"
	supplierDisplayName := "test supplier"
	consumerOperationResourceName := "test"
	consumerExternalApiResourceName := ""
	supplierOperationResourceName := "test"
	supplierExternalApiResourceName := ""

	if _, err = CreateDependency(dependencyID, description, consumerDisplayName,
		consumerOperationResourceName, consumerExternalApiResourceName, supplierDisplayName,
		supplierOperationResourceName, supplierExternalApiResourceName); err != nil {
		t.Errorf("failed to create dependency %v", err)
	}
}

func TestGetDependency(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	dependencyID := "test-dependency"

	if _, err := GetDependency(dependencyID); err != nil {
		t.Errorf("failed to get dependency %v", err)
	}
}

func TestListDependencies(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	if _, err := ListDependencies("", -1, ""); err != nil {
		t.Errorf("failed to list dependencies %v", err)
	}
}

func TestDeleteDependency(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	dependencyID := "test-dependency"

	if _, err := DeleteDependency(dependencyID); err != nil {
		t.Errorf("failed to get dependency %v", err)
	}
}

func TestCreateDeployment(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	deploymentID := "test-deployment"
	description := "test description"
	displayName := "test display name"
	externalURI := "https://httpbin.org/get"
	resourceURI := "https://httpbin.org/get"
	endpoints := []string{"https://httpbin.org/get"}

	if _, err = CreateDeployment(deploymentID, displayName, description, externalURI,
		resourceURI, endpoints, APIGEE, DEVELOPMENT, SLO99_9); err != nil {
		t.Errorf("failed to create deployment %v", err)
	}
}

func TestGetDeployment(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	deploymentID := "test-deployment"

	if _, err := GetDeployment(deploymentID); err != nil {
		t.Errorf("failed to get deployment %v", err)
	}
}

func TestListDeployment(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	if _, err := ListDeployments("", -1, ""); err != nil {
		t.Errorf("failed to list deployments %v", err)
	}
}

func TestDeleteDeployment(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	deploymentID := "test-deployment"

	if _, err := DeleteDeployment(deploymentID); err != nil {
		t.Errorf("failed to delete deployment %v", err)
	}
}

func TestCreateExternalApis(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	externalApiId := "test-external-api"
	description := "test description"
	displayName := "test display name"
	externalURI := "https://httpbin.org/get"
	paths := []string{"/get"}
	endpoints := []string{"https://httpbin.org/get"}

	if _, err = CreateExternalAPI(externalApiId, displayName, description, endpoints, paths, externalURI); err != nil {
		t.Errorf("failed to create external api %v", err)
	}
}

func TestGetExternalApi(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	externalApiId := "test-external-api"

	if _, err := GetExternalAPI(externalApiId); err != nil {
		t.Errorf("failed to get external api %v", err)
	}
}

func TestGetInstance(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	instanceId := "test-instance"

	if _, err := GetInstance(instanceId); err != nil {
		t.Errorf("failed to get instance %v", err)
	}
}

func TestLookupInstance(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if os.Getenv("APIGEE_REGION") == "" {
		t.Fatalf("APIGEE_REGION not set")
	}

	if _, err := LookupInstance(); err != nil {
		t.Errorf("failed to lookup instances %v", err)
	}
}

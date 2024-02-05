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

package apis

import (
	"testing"

	"internal/apiclient"
	"internal/client/clienttest"
)

func TestCreate(t *testing.T) {
	name := "test"
	displayName := "test"
	description := "description"
	recommendedVersion := "projects/my-project/locations/us-west1/apis/test/versions/1"
	recommendedDeployment := "projects/my-project/locations/us-west1/apis/test/deployments/1"
	labels := map[string]string{
		"key1": "value1",
	}
	annotations := map[string]string{
		"key1": "value1",
	}

	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if apiclient.GetRegion() == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	if _, err := Create("id123", name, displayName, description,
		"NONE", recommendedVersion, recommendedDeployment, labels, annotations); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGet(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if apiclient.GetRegion() == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	if _, err := Get("id123"); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestList(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if apiclient.GetRegion() == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	if _, err := List(-1, "", "", ""); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if apiclient.GetRegion() == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	if _, err := Delete("id123"); err != nil {
		t.Fatalf("%v", err)
	}
}

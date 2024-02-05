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

package artifacts

import (
	"encoding/base64"
	"os"
	"path"
	"testing"

	"internal/apiclient"
	"internal/client/clienttest"
	"internal/cmd/utils"
)

var cliPath = os.Getenv("APIGEECLI_PATH")

func TestCreate(t *testing.T) {
	artifactID := "test"
	name := "test"
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if apiclient.GetRegion() == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	payload, err := utils.ReadFile(path.Join(cliPath, "test", "openapi.json"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	contents := base64.StdEncoding.EncodeToString(payload)
	labels := map[string]string{
		"test": "test",
	}
	annotations := map[string]string{
		"test": "test",
	}
	if _, err := Create(artifactID, name, string(contents), labels, annotations); err != nil {
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
	if _, err := Get("test"); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetContents(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if apiclient.GetRegion() == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	if err := GetContents("test"); err != nil {
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

func TestReplace(t *testing.T) {
	artifactID := "test"
	name := "test"
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if apiclient.GetRegion() == "" {
		t.Fatalf("APIGEE_REGION not set")
	}
	payload, err := utils.ReadFile(path.Join(cliPath, "test", "openapi.json"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	contents := base64.StdEncoding.EncodeToString(payload)
	labels := map[string]string{
		"test": "test",
	}
	annotations := map[string]string{
		"test": "test",
	}
	if _, err := Replace(artifactID, name, string(contents), labels, annotations); err != nil {
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
	if _, err := Delete("test"); err != nil {
		t.Fatalf("%v", err)
	}
}

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

package apps

import (
	"encoding/json"
	"fmt"
	"internal/client/clienttest"
	"internal/client/developers"
	"internal/client/products"
	"os"
	"testing"
)

const (
	email = "user@example.com"
	name  = "test"
)

var appID, devID string

func TestCreate(t *testing.T) {
	var respBody []byte
	var respJSONMap map[string]interface{}
	var err error
	if err = clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	fName := "Test"
	lName := "User"
	username := "testuser"
	attrs := map[string]string{
		"test": "test",
	}
	if respBody, err = developers.Create(email, fName, lName, username, attrs); err != nil {
		t.Fatalf("%v", err)
	}
	if err = json.Unmarshal(respBody, &respJSONMap); err != nil {
		t.Fatalf("%v", err)
	}
	devID = respJSONMap["id"].(string)

	createProduct(t)

	expires := "-1"
	callback := ""
	scopes := []string{"test"}
	apiProducts := []string{"test"}

	if respBody, err = Create(name, email, expires, callback, apiProducts, scopes, attrs); err != nil {
		t.Fatalf("%v", err)
	}
	if err = json.Unmarshal(respBody, &respJSONMap); err != nil {
		t.Fatalf("%v", err)
	}
	appID = respJSONMap["id"].(string)
	if appID == "" {
		t.Fatalf("%v", fmt.Errorf("unable to find appId"))
	}
}

func TestGet(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Get(appID); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestSearch(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := SearchApp(name); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestManage(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Manage(appID, email, "approve"); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestList(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := List(false, true, -1, "", "", "", "", "", -1, "", ""); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListApps(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListApps("test"); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGenerateKey(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	apiProducts := []string{"test"}
	expires := "-1"
	callback := ""
	scopes := []string{"test"}
	if _, err := GenerateKey(name, devID, apiProducts, callback, expires, scopes); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestExport(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Export(4); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Delete(name, devID); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := developers.Delete(email); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := products.Delete("test"); err != nil {
		t.Fatalf("%v", err)
	}
}

func createProduct(t *testing.T) {
	p := products.APIProduct{}

	p.Name = "test"
	p.DisplayName = "test"
	p.ApprovalType = ""
	p.Description = ""
	p.Environments = []string{os.Getenv("APIGEE_ENV")}

	if _, err := products.Create(p); err != nil {
		t.Fatalf("%v", err)
	}
}

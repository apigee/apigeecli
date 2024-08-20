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

package appgroups

import (
	"internal/client/clienttest"
	"internal/client/products"
	"os"
	"testing"
)

const (
	appName = "test"
)

func TestCreateApp(t *testing.T) {
	expires := "-1"
	callback := ""
	scopes := []string{"test"}
	apiProducts := []string{"test"}
	attrs := map[string]string{
		"test": "test",
	}
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	createProduct(t)

	TestCreate(t)

	if _, err := CreateApp(name, appName, expires, callback, apiProducts, scopes, attrs); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetApp(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetApp(name, appName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListApps(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListApps(name, -1, ""); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestManageApp(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ManageApp(name, appName, "approve"); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUpdateApp(t *testing.T) {
	expires := ""
	callback := ""
	scopes := []string{"update"}
	apiProducts := []string{"test"}
	attrs := map[string]string{
		"test": "update",
	}
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	if _, err := UpdateApp(name, appName, expires, callback, apiProducts, scopes, attrs); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestExportApps(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ExportApps(name); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeleteApp(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteApp(name, appName); err != nil {
		t.Fatalf("%v", err)
	}
	TestDelete(t)
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

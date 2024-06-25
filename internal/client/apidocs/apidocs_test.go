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

package apidocs

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"internal/client/clienttest"
	"internal/cmd/utils"
)

var (
	aPIDocID = ""
	title    = "test"
	siteID   = os.Getenv("APIGEE_SITEID")
	cliPath  = os.Getenv("APIGEECLI_PATH")
)

func TestCreate(t *testing.T) {
	var respBody []byte
	var respJSONMap map[string]interface{}
	description := "description for test"
	published := "false"
	anonAllowed := "false"
	requireCallbackUrl := "false"
	imageUrl := ""
	categoryIds := []string{}
	apiProductName := "test"

	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if respBody, err = Create(siteID, title, description, published, anonAllowed,
		apiProductName, requireCallbackUrl, imageUrl, categoryIds); err != nil {
		t.Fatalf("%v", err)
	}
	if err = json.Unmarshal(respBody, &respJSONMap); err != nil {
		t.Fatalf("%v", err)
	}
	if data, ok := respJSONMap["data"].(map[string]interface{}); ok {
		if data["id"].(string) != "" {
			aPIDocID = data["id"].(string)
		} else {
			t.Fatalf("%v", fmt.Errorf("ID not found or not a string"))
		}
	} else {
		t.Fatalf("%v", fmt.Errorf("Data not found or not a map"))
	}
}

func TestGet(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if _, err := Get(siteID, aPIDocID); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetByTitle(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if _, err := GetByTitle(siteID, title); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestList(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if _, err := List(siteID, -1, ""); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUpdateDocumentation(t *testing.T) {
	var openAPIDoc []byte
	displayName := "test"
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if openAPIDoc, err = utils.ReadFile(path.Join(cliPath, "test", "openapi.json")); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err = UpdateDocumentation(siteID, aPIDocID, displayName, openAPIDoc, nil, ""); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetDocumentation(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if _, err := GetDocumentation(siteID, aPIDocID); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if _, err := Delete(siteID, aPIDocID); err != nil {
		t.Fatalf("%v", err)
	}
}

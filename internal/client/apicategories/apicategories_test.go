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

package apicategories

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"internal/client/clienttest"
)

const (
	categoryName = "test"
)

var (
	siteID        = os.Getenv("APIGEE_SITEID")
	aPICategoryID string
)

func TestCreate(t *testing.T) {
	var respBody []byte
	var respJSONMap map[string]interface{}
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if respBody, err = Create(siteID, categoryName); err != nil {
		t.Fatalf("%v", err)
	}

	if err = json.Unmarshal(respBody, &respJSONMap); err != nil {
		t.Fatalf("%v", err)
	}
	if data, ok := respJSONMap["data"].(map[string]interface{}); ok {
		if data["id"].(string) != "" {
			aPICategoryID = data["id"].(string)
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
	if _, err := Get(siteID, aPICategoryID); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetByName(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if _, err := GetByName(siteID, categoryName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetIDByName(t *testing.T) {
	var id string
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if id, err = GetIDByName(siteID, categoryName); err != nil {
		t.Fatalf("%v", err)
	}
	if id != aPICategoryID {
		t.Fatalf("IDs do not match: %s != %s", id, aPICategoryID)
	}
}

func TestList(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if _, err := List(siteID); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	if _, err := Delete(siteID, aPICategoryID); err != nil {
		t.Fatalf("%v", err)
	}
}

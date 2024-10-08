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
	"internal/client/clienttest"
	"testing"
)

func TestCreateKey(t *testing.T) {
	scopes := []string{"test"}
	apiProducts := []string{"test"}
	attrs := map[string]string{
		"test": "test",
	}
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	TestCreate(t)

	if _, err := CreateKey(name, appID, "key1", "key1-secret", apiProducts, scopes, "-1", attrs); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetKey(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetKey(name, appID, "key1"); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestManageKey(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ManageKey(name, appID, "key1", "approve"); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeleteKey(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteKey(name, appID, "key1"); err != nil {
		t.Fatalf("%v", err)
	}
	TestDelete(t)
}

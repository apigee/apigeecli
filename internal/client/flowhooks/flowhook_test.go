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

package flowhooks

import (
	"internal/client/clienttest"
	"internal/client/sharedflows"
	"os"
	"path"
	"testing"
)

var cliPath = os.Getenv("APIGEECLI_PATH")

const (
	name       = "test-sharedflow"
	testFolder = "test"
)

func TestAttach(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := sharedflows.Create(name, path.Join(cliPath, testFolder, "test_flow.zip")); err != nil {
		t.Fatalf("%v", err)
	}
	cPtr := new(bool)
	*cPtr = true
	if _, err := Attach("PreProxyFlowHook", "test description", name, cPtr); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGet(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Get(name); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestList(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := List(); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Detach(name); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := sharedflows.Delete(name, -1); err != nil {
		t.Fatalf("%v", err)
	}
}

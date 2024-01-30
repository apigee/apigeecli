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

package sharedflows

import (
	"internal/client/clienttest"
	"os"
	"path"
	"testing"
)

const (
	sfName     = "test-sf"
	testFolder = "test"
)

var cliPath = os.Getenv("APIGEECLI_PATH")

func TestCreate(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Create(sfName, path.Join(cliPath, testFolder, "test_flow.zip")); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGet(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Get(sfName, 1); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Get(sfName, -1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetHighestSfRevision(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetHighestSfRevision(sfName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestList(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := List(true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeploy(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Deploy(sfName, 1, false, ""); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListDeployments(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	if _, err := ListDeployments(sfName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListRevisionDeployments(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	if _, err := ListRevisionDeployments(sfName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListEnvDeployments(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	if _, err := ListEnvDeployments(); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUndeploy(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Undeploy(sfName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestFetch(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if err := Fetch(sfName, 1); err != nil {
		t.Fatalf("%v", err)
	}
	os.Remove(sfName + ".zip")
}

func TestDelete(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Delete(sfName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

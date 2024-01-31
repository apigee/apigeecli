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

package env

import (
	"internal/apiclient"
	"internal/client/clienttest"
	"internal/cmd/utils"
	"os"
	"path"
	"testing"
)

var cliPath = os.Getenv("APIGEECLI_PATH")

func TestSetDebug(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}

	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv(name)

	TestCreate(t)

	payload, err := utils.ReadFile(path.Join(cliPath, "test", "debugmask.json"))
	if err != nil {
		t.Fatalf("%v", err)
	}

	if _, err := SetDebug(payload); err != nil {
		t.Fatalf("%v", err)
	}
	apiclient.SetApigeeEnv(env)
}

func TestGetDebug(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}

	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv(name)
	if _, err := GetDebug(); err != nil {
		t.Fatalf("%v", err)
	}

	TestDelete(t)
	apiclient.SetApigeeEnv(env)
}

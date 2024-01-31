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
	"testing"
)

const name = "unittest-env"

func TestCreate(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}

	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv(name)

	if _, err := Create("BASIC", "PROXY", ""); err != nil {
		t.Fatalf("%v", err)
	}

	apiclient.SetApigeeEnv(env)
}

func TestGet(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}
	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv(name)

	if _, err := Get(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Get(true); err != nil {
		t.Fatalf("%v", err)
	}

	apiclient.SetApigeeEnv(env)
}

func TestList(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if _, err := List(); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeployments(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if _, err := GetDeployments(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetAllDeployments(); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetDeployedConfig(); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestSetProperty(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}
	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv(name)
	if err := SetEnvProperty("test", "test"); err != nil {
		t.Fatalf("%v", err)
	}
	apiclient.SetApigeeEnv(env)
}

func TestClearEnvProperties(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}
	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv(name)
	if err := ClearEnvProperties(); err != nil {
		t.Fatalf("%v", err)
	}
	apiclient.SetApigeeEnv(env)
}

func TestExport(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Export(); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD)
	if err != nil {
		t.Fatalf("%v", err)
	}
	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv(name)

	if _, err := Delete(); err != nil {
		t.Fatalf("%v", err)
	}

	apiclient.SetApigeeEnv(env)
}

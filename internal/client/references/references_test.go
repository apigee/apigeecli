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

package references

import (
	"internal/apiclient"
	"internal/client/clienttest"
	"internal/client/env"
	"internal/client/keystores"
	"testing"
)

const (
	keyStoreName = "test"
	name         = "test"
)

func TestCreate(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	apiclient.SetApigeeEnv(name)
	if _, err := env.Create("BASIC", "PROXY", ""); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := keystores.Create(keyStoreName); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Create(name, "description", "KeyStore", keyStoreName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGet(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	apiclient.SetApigeeEnv(name)
	if _, err := Get(name); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestList(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	apiclient.SetApigeeEnv(name)
	if _, err := List(); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUpdate(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	apiclient.SetApigeeEnv(name)
	if _, err := Update(name, "change description", "KeyStore", keyStoreName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	apiclient.SetApigeeEnv(name)
	if _, err := Delete(name); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := keystores.Delete(keyStoreName); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := env.Delete(); err != nil {
		t.Fatalf("%v", err)
	}
}

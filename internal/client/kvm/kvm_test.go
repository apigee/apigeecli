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

package kvm

import (
	"os"
	"path"
	"testing"

	"internal/apiclient"
	"internal/client/apis"
	"internal/client/clienttest"
)

const (
	proxyName  = "test"
	testFolder = "test"
	kvmName    = "test"
)

var cliPath = os.Getenv("APIGEECLI_PATH")

func TestCreate(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv("")

	// test org KVM
	if _, err := Create("", kvmName, true); err != nil {
		t.Fatalf("%v", err)
	}

	// test proxy KVM
	if _, err := apis.CreateProxy(proxyName, path.Join(cliPath, testFolder, "test_proxy.zip")); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Create(proxyName, kvmName, true); err != nil {
		t.Fatalf("%v", err)
	}

	// test env KVM
	apiclient.SetApigeeEnv(env)
	if _, err := Create("", kvmName, true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestList(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv("")

	// test org KVM
	if _, err := Delete("", kvmName); err != nil {
		t.Fatalf("%v", err)
	}

	// test proxy KVM
	if _, err := Delete(proxyName, kvmName); err != nil {
		t.Fatalf("%v", err)
	}

	// test env KVM
	apiclient.SetApigeeEnv(env)
	if _, err := Delete("", kvmName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCleanup(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := apis.DeleteProxy(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

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
	"testing"

	"internal/apiclient"
	"internal/client/clienttest"
)

const (
	keyName = "test"
	value   = "test"
)

func TestCreateEntry(t *testing.T) {
	TestCreate(t)

	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv("")

	// add entry to org kvm
	if _, err := CreateEntry("", kvmName, keyName, value); err != nil {
		t.Fatalf("%v", err)
	}

	// add entry to proxy kvm
	if _, err := CreateEntry(proxyName, kvmName, keyName, value); err != nil {
		t.Fatalf("%v", err)
	}

	// add entry to env kvm
	apiclient.SetApigeeEnv(env)
	if _, err := CreateEntry("", kvmName, keyName, value); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetEntry(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv("")

	// test org kvm entry
	if _, err := GetEntry("", kvmName, keyName); err != nil {
		t.Fatalf("%v", err)
	}

	// test proxy kvm entry
	if _, err := GetEntry(proxyName, kvmName, keyName); err != nil {
		t.Fatalf("%v", err)
	}

	// test env kvm entry
	apiclient.SetApigeeEnv(env)
	if _, err := GetEntry("", kvmName, keyName); err != nil {
		t.Fatalf("%v", err)
	}

	// TODO: assert values
}

func TestUpdateEntry(t *testing.T) {
	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv("")

	updatedValue := "updatedTest"

	// add entry to org kvm
	if _, err := UpdateEntry("", kvmName, keyName, updatedValue); err != nil {
		t.Fatalf("%v", err)
	}

	// add entry to proxy kvm
	if _, err := UpdateEntry(proxyName, kvmName, keyName, updatedValue); err != nil {
		t.Fatalf("%v", err)
	}

	// add entry to env kvm
	apiclient.SetApigeeEnv(env)
	if _, err := UpdateEntry("", kvmName, keyName, updatedValue); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListEntries(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv("")

	// test org kvm entry
	if _, err := ListEntries("", kvmName, -1, ""); err != nil {
		t.Fatalf("%v", err)
	}

	// test proxy kvm entry
	if _, err := ListEntries(proxyName, kvmName, -1, ""); err != nil {
		t.Fatalf("%v", err)
	}

	// test env kvm entry
	apiclient.SetApigeeEnv(env)
	if _, err := ListEntries("", kvmName, -1, ""); err != nil {
		t.Fatalf("%v", err)
	}

	// TODO: assert values
}

func TestDeleteEntry(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	env := apiclient.GetApigeeEnv()
	apiclient.SetApigeeEnv("")

	// test org kvm entry
	if _, err := DeleteEntry("", kvmName, keyName); err != nil {
		t.Fatalf("%v", err)
	}

	// test proxy kvm entry
	if _, err := DeleteEntry(proxyName, kvmName, keyName); err != nil {
		t.Fatalf("%v", err)
	}

	// test env kvm entry
	apiclient.SetApigeeEnv(env)
	if _, err := DeleteEntry("", kvmName, keyName); err != nil {
		t.Fatalf("%v", err)
	}

	// TODO: assert values
}

func TestEntriesCleanup(t *testing.T) {
	TestDelete(t)
	TestCleanup(t)
}

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

package appgroups

import (
	"testing"

	"internal/client/clienttest"
)

const (
	name = "test"
)

func TestCreate(t *testing.T) {
	channelURI := "https://example.com/channel"
	channelID := "test"
	displayName := "test"
	attrs := map[string]string{
		"test": "test",
	}
	devs := map[string]string{
		"test": "test",
	}
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Create(name, channelURI, channelID, displayName, attrs, devs); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGet(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Get(name); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestList(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := List(-1, "", ""); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestManage(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Manage(name, "active"); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUpdate(t *testing.T) {
	channelURI := "https://example.com/update"
	channelID := ""
	displayName := ""
	attrs := map[string]string{
		"test": "update",
	}
	devs := map[string]string{
		"test": "update",
	}
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Update(name, channelURI, channelID, displayName, attrs, devs); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestExport(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Export(); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Delete(name); err != nil {
		t.Fatalf("%v", err)
	}
}

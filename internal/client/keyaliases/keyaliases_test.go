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

package keyaliases

import (
	"internal/client/clienttest"
	"internal/client/keystores"
	"os"
	"path"
	"testing"
)

const (
	keyStoreName = "test"
	name         = "test"
)

var cliPath = os.Getenv("APIGEECLI_PATH")

func TestCreateKeystore(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := keystores.Create(keyStoreName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCreateOrUpdateSelfSigned(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := CreateOrUpdateSelfSigned(keyStoreName, "test", false, true, true,
		path.Join(cliPath, "test", "self-signed.json")); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Delete(keyStoreName, name); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCreateOrUpdatePfx(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	password := os.Getenv("APIGEE_KEYALIAS_PWD")
	if password == "" {
		t.Fatalf("APIGEE_KEYALIAS_PWD not set")
	}
	if _, err := CreateOrUpdatePfx(keyStoreName, "test", false, true, true,
		path.Join(cliPath, "test", "pkcs12.pfx"), password); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Delete(keyStoreName, name); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCreateOrUpdateKeyCert(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	password := os.Getenv("APIGEE_KEYALIAS_PWD")
	if password == "" {
		t.Fatalf("APIGEE_KEYALIAS_PWD not set")
	}
	if _, err := CreateOrUpdateKeyCert(keyStoreName, "test", false, true, true,
		path.Join(cliPath, "test", "cert.pem"), path.Join(cliPath, "test", "key.pem"), password); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetCert(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if err := GetCert(keyStoreName, "test"); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Delete(keyStoreName, name); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCreateCSR(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := CreateCSR(keyStoreName, "test"); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Delete(keyStoreName, name); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestList(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := List(keyStoreName); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := keystores.Delete(keyStoreName); err != nil {
		t.Fatalf("%v", err)
	}
}

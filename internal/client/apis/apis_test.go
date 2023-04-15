// Copyright 2022 Google LLC
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

package apis

import (
	"fmt"
	"os"
	"path"
	"testing"

	"internal/apiclient"
)

const (
	proxyName  = "test"
	testFolder = "test"
)

func setup() (err error) {
	org := os.Getenv("APIGEE_ORG")
	if org == "" {
		return fmt.Errorf("APIGEE_ORG not set")
	}

	if err = apiclient.SetApigeeOrg(org); err != nil {
		return fmt.Errorf("APIGEE_ORG not set")
	}

	token := os.Getenv("APIGEE_TOKEN")
	if token == "" {
		return fmt.Errorf("APIGEE_TOKEN not set")
	}
	apiclient.SetApigeeToken(token)
	return nil
}

func getApigeecliHome() (cliPath string, err error) {
	cliPath = os.Getenv("APIGEECLI_PATH")
	if cliPath == "" {
		return "", fmt.Errorf("APIGEECLI_PATH not set")
	}
	return cliPath, err
}

func TestCreateProxy(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("%v", err)
	}
	cliPath, err := getApigeecliHome()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if _, err = CreateProxy(proxyName, path.Join(cliPath, testFolder, "test_proxy.zip")); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeleteProxy(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteProxy(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeleteProxyRevision(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteProxyRevision(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestFetchProxy(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("%v", err)
	}
	if err := FetchProxy(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetProxy(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetProxy(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetProxy(proxyName, -1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetHighestProxyRevision(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetHighestProxyRevision(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

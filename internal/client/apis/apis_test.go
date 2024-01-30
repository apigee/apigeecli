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
	"internal/client/clienttest"
	"os"
	"path"
	"testing"
)

// go test internal/client/apis

const (
	proxyName  = "test"
	testFolder = "test"
)

var cliPath = os.Getenv("APIGEECLI_PATH")

func TestCreateProxy(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_REQD); err != nil {
		t.Fatalf("%v", err)
	}

	if _, err := CreateProxy(proxyName, path.Join(cliPath, testFolder, "test_proxy.zip")); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestFetchProxy(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if err := FetchProxy(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
	os.Remove(proxyName + "1.zip")
}

func TestGetProxy(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
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
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetHighestProxyRevision(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGenerateDeployChangeReport(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GenerateDeployChangeReport(proxyName, 1, false); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGenerateUndeployChangeReport(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GenerateUndeployChangeReport(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListProxies(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListProxies(true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListProxyDeployments(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListProxyDeployments(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListProxyRevisionDeployments(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListProxyRevisionDeployments(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeployProxy(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeployProxy(proxyName, 1, false, false, false, ""); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUndeployProxy(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := UndeployProxy(proxyName, 1, true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUpdate(t *testing.T) {
	labels := map[string]string{
		"label1": "value1",
		"label2": "value2",
	}
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Update(proxyName, labels); err != nil {
		t.Fatalf("%v", err)
	}
}

// TODO: CleanProxy, ExportProxies

func TestDeleteProxyRevision(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := CreateProxy(proxyName, ""); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteProxyRevision(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeleteProxy(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteProxy(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

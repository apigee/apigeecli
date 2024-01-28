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

// go test internal/client/apis

const (
	proxyName  = "test"
	testFolder = "test"
)

func setup(envReqd bool) (err error) {
	org := os.Getenv("APIGEE_ORG")
	if org == "" {
		return fmt.Errorf("APIGEE_ORG not set")
	}

	if envReqd {
		env := os.Getenv("APIGEE_ENV")
		if env == "" {
			return fmt.Errorf("APIGEE_ENV not set")
		}
	}

	apiclient.NewApigeeClient(apiclient.ApigeeClientOptions{
		TokenCheck:  true,
		PrintOutput: true,
		NoOutput:    false,
		DebugLog:    true,
		SkipCache:   false,
	})

	if err = apiclient.SetApigeeOrg(org); err != nil {
		return fmt.Errorf("APIGEE_ORG not set")
	}

	if apiclient.GetApigeeToken() == "" {
		token := os.Getenv("APIGEE_TOKEN")
		if token == "" {
			err = apiclient.GetDefaultAccessToken()
			if err != nil {
				return fmt.Errorf("APIGEE_TOKEN not set")
			}
		} else {
			apiclient.SetApigeeToken(token)
		}
	}

	_, err = getApigeecliHome()
	if err != nil {
		return err
	}

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
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	cliPath, _ := getApigeecliHome()
	if _, err := CreateProxy(proxyName, path.Join(cliPath, testFolder, "test_proxy.zip")); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestFetchProxy(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if err := FetchProxy(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
	os.Remove(proxyName + "1.zip")
}

func TestGetProxy(t *testing.T) {
	if err := setup(false); err != nil {
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
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GetHighestProxyRevision(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGenerateDeployChangeReport(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GenerateDeployChangeReport(proxyName, 1, false); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGenerateUndeployChangeReport(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := GenerateUndeployChangeReport(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListProxies(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListProxies(true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListProxyDeployments(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListProxyDeployments(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListProxyRevisionDeployments(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListProxyRevisionDeployments(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeployProxy(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeployProxy(proxyName, 1, false, false, false, ""); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUndeployProxy(t *testing.T) {
	if err := setup(false); err != nil {
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
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Update(proxyName, labels); err != nil {
		t.Fatalf("%v", err)
	}
}

// TODO: CleanProxy, ExportProxies

func TestDeleteProxyRevision(t *testing.T) {
	if err := setup(false); err != nil {
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
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteProxy(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

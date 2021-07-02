// Copyright 2020 Google LLC
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

package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
)

const apigeecli = "./apigeecli"

var token = os.Getenv("APIGEE_TOKEN")
var env = os.Getenv("APIGEE_ENV")
var org = os.Getenv("APIGEE_ORG")

const email = "test@gmail.com"
const name = "test"
const attrs = "\"foo1=bar1,foo2=bar2\""
const envs = "test"

var appID string

func TestMain(t *testing.T) {
	cmd := exec.Command(apigeecli)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// orgs test
func TestListOrgs(t *testing.T) {
	cmd := exec.Command(apigeecli, "orgs", "list", "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOrg(t *testing.T) {
	cmd := exec.Command(apigeecli, "orgs", "get", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// env tests
func TestListEnvs(t *testing.T) {
	cmd := exec.Command(apigeecli, "envs", "list", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetEnv(t *testing.T) {
	cmd := exec.Command(apigeecli, "envs", "get", "-o", org, "-e", env, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// developers test
func TestCreateDeveloper(t *testing.T) {
	first := "frstname"
	last := "lastname"
	user := "username"

	cmd := exec.Command(apigeecli, "developers", "create", "-o", org, "-n", email, "-f", first, "-s",
		last, "-u", user, "--attrs", attrs, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDeveloper(t *testing.T) {

	cmd := exec.Command(apigeecli, "developers", "get", "-o", org, "-n", email, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListDeveloper(t *testing.T) {

	cmd := exec.Command(apigeecli, "developers", "list", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListExpandDeveloper(t *testing.T) {
	expand := "true"

	cmd := exec.Command(apigeecli, "developers", "list", "-o", org, "-x", expand, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// kvm test
func TestCreateKVM(t *testing.T) {

	cmd := exec.Command(apigeecli, "kvms", "create", "-o", org, "-e", env, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListKVM(t *testing.T) {

	cmd := exec.Command(apigeecli, "kvms", "list", "-o", org, "-e", env, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// sync tests
func TestGetSync(t *testing.T) {

	cmd := exec.Command(apigeecli, "sync", "get", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetSync(t *testing.T) {
	ity := "test@gmail.com"
	cmd := exec.Command(apigeecli, "sync", "set", "-o", org, "-i", ity, "-t", token)
	err := cmd.Run()
	if err == nil {
		t.Fatal("Invalid identity test should have failed")
	}
}

func TestCreateProxy(t *testing.T) {
	cmd := exec.Command(apigeecli, "apis", "create", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateProductLegacy(t *testing.T) {
	proxy := "test"
	approval := "auto"

	cmd := exec.Command(apigeecli, "products", "create", "-o", org, "-n", name+"-legacy", "-e", envs, "-p", proxy,
		"-f", approval, "--attrs", attrs, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateProductGroup(t *testing.T) {
	approval := "auto"
	operationGroupFile := "./test/apiproduct-op-group.json"

	cmd := exec.Command(apigeecli, "products", "create", "-o", org, "-n", name+"-group", "-e", envs,
		"-g", operationGroupFile, "-f", approval, "--attrs", attrs, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateApp(t *testing.T) {
	product := "test"

	cmd := exec.Command(apigeecli, "apps", "create", "-o", org, "-n", name, "-e", email, "-p", product,
		"--attrs", attrs, "-t", token)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	var outputMap map[string]string
	_ = json.Unmarshal(out.Bytes(), &outputMap)
	appID = outputMap["appId"]
}

func TestListApp(t *testing.T) {

	cmd := exec.Command(apigeecli, "apps", "list", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetProduct(t *testing.T) {

	cmd := exec.Command(apigeecli, "products", "get", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListProduct(t *testing.T) {

	cmd := exec.Command(apigeecli, "products", "list", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDeveloper(t *testing.T) {

	cmd := exec.Command(apigeecli, "developers", "delete", "-o", org, "-n", email, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateTS(t *testing.T) {
	host := "example.com"

	cmd := exec.Command(apigeecli, "targetservers", "create", "-o", org, "-n", name, "-e", env, "-s",
		host, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetTS(t *testing.T) {

	cmd := exec.Command(apigeecli, "targetservers", "get", "-o", org, "-e", env, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestListTS(t *testing.T) {

	cmd := exec.Command(apigeecli, "targetservers", "list", "-o", org, "-e", env, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestDeleteTS(t *testing.T) {

	cmd := exec.Command(apigeecli, "targetservers", "delete", "-o", org, "-e", env, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestExportProducts(t *testing.T) {

	cmd := exec.Command(apigeecli, "products", "export", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteKVM(t *testing.T) {

	cmd := exec.Command(apigeecli, "kvms", "delete", "-o", org, "-e", env, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteProduct(t *testing.T) {

	cmd := exec.Command(apigeecli, "products", "delete", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateEncKVM(t *testing.T) {
	name := "testEnc"

	cmd := exec.Command(apigeecli, "kvms", "create", "-o", org, "-e", env, "-n", name, "-c", "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEncKVM(t *testing.T) {
	name := "testEnc"

	cmd := exec.Command(apigeecli, "kvms", "delete", "-o", org, "-e", env, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteProxy(t *testing.T) {
	cmd := exec.Command(apigeecli, "apis", "delete", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateResource(t *testing.T) {
	name := "test.js"
	resType := "jsc"
	resPath := "./test/test.js"
	cmd := exec.Command(apigeecli, "resources", "create", "-o", org, "-e", env, "-n", name, "-p", resType, "-r", resPath)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetResource(t *testing.T) {
	name := "test.js"
	resType := "jsc"
	cmd := exec.Command(apigeecli, "resources", "get", "-o", org, "-e", env, "-n", name, "-p", resType)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListResources(t *testing.T) {
	cmd := exec.Command(apigeecli, "resources", "list", "-o", org, "-e", env)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteResource(t *testing.T) {
	name := "test.js"
	resType := "jsc"
	cmd := exec.Command(apigeecli, "resources", "delete", "-o", org, "-e", env, "-n", name, "-p", resType)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

/*func TestCreateKeyStore(t *testing.T) {
	name := "testkeystore"
	cmd := exec.Command(apigeecli, "keystores", "create", "-o", org, "-e", env,
		"-n", name)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateKeyAliasSelfSigned(t *testing.T) {
	payload := "{\"alias\": \"self\",\"keySize\": \"2048\",\"sigAlg\": \"SHA256withRSA\",\"subject\"" +
		"{\"commonName\": \"nandan\",\"email\": \"test@test.com\"},\"certValidityInDays\": \"10\"}"
	name := "testkeystore"
	keyaliasname := "selfkeyalias"
	cmd := exec.Command(apigeecli, "keyaliases", "create-self-signed", "-o", org, "-e", env,
		"-n", name, "-s", keyaliasname, "-c", payload)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}*/

//get requires an app id. run list and get an app id or get it from create
/*func TestGetApp(t *testing.T) {
	cmd := exec.Command(apigeecli, "apps", "get", "-o", org, "-i", appID, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}*/

/*func TestSetMart(t *testing.T) {
	mart := os.Getenv("MART")
	if mart == "" {
		t.Log("MART not set, skipping")
	} else {
		cmd := exec.Command(apigeecli, "orgs", "setmart", "-o", org, "-m", mart, "-t", token)
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSetMartWhiteList(t *testing.T) {
	mart := os.Getenv("MART")
	if mart == "" {
		t.Log("MART not set, skipping")
	} else {
		cmd := exec.Command(apigeecli, "orgs", "setmart", "-o", org, "-m", mart, "-w", "false", "-t", token)
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeleteApp(t *testing.T) {

	cmd := exec.Command(apigeecli, "apps", "delete", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}
*/

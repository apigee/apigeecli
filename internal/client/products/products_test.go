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

package products

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"internal/client/clienttest"
	"internal/clilog"
)

const name = "test"

var cliPath = os.Getenv("APIGEECLI_PATH")

func TestCreate(t *testing.T) {
	var err error
	if err = clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	p := APIProduct{}

	p.Name = name
	p.DisplayName = "test"
	p.ApprovalType = ""
	p.Description = ""
	p.Environments = []string{os.Getenv("APIGEE_ENV")}
	p.OperationGroup, err = getOperationGroup(path.Join(cliPath, "test", "apiproduct-op-group.json"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	p.GraphQLOperationGroup, err = getGqlOperationGroup(path.Join(cliPath, "test", "apiproduct-gqlop-group.json"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	p.GrpcOperationGroup, err = getGrpcOperationGroup(path.Join(cliPath, "test", "apiproduct-grpcop-group.json"))
	if err != nil {
		t.Fatalf("%v", err)
	}

	if _, err = Create(p); err != nil {
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
	if _, err := List(-1, "", true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestExport(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := Export(4); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUpdate(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	p := APIProduct{}

	p.Name = name
	p.Description = "adding description"
	if _, err := Update(p); err != nil {
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

func getOperationGroup(operationGroupFile string) (*OperationGroup, error) {
	var operationGrpBytes []byte
	var err error

	operationGrp := OperationGroup{}

	if operationGroupFile != "" {
		if operationGrpBytes, err = os.ReadFile(operationGroupFile); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(operationGrpBytes, &operationGrp); err != nil {
			return nil, err
		}
		return &operationGrp, nil
	}
	return nil, nil
}

func getGrpcOperationGroup(grpcOperationGroupFile string) (*GrpcOperationGroup, error) {
	var grpcOperationGrpBytes []byte
	var err error

	grpcOperationGrp := GrpcOperationGroup{}

	if grpcOperationGroupFile != "" {
		if grpcOperationGrpBytes, err = os.ReadFile(grpcOperationGroupFile); err != nil {
			clilog.Debug.Println(err)
			return nil, err
		}
		if err = json.Unmarshal(grpcOperationGrpBytes, &grpcOperationGrp); err != nil {
			clilog.Debug.Println(err)
			return nil, err
		}
		return &grpcOperationGrp, nil
	}
	return nil, nil
}

func getGqlOperationGroup(gqlOperationGroupFile string) (*GraphqlOperationGroup, error) {
	var gqlOperationGrpBytes []byte
	var err error

	gqlOperationGrp := GraphqlOperationGroup{}

	if gqlOperationGroupFile != "" {
		if gqlOperationGrpBytes, err = os.ReadFile(gqlOperationGroupFile); err != nil {
			clilog.Debug.Println(err)
			return nil, err
		}
		if err = json.Unmarshal(gqlOperationGrpBytes, &gqlOperationGrp); err != nil {
			clilog.Debug.Println(err)
			return nil, err
		}
		return &gqlOperationGrp, nil
	}
	return nil, nil
}

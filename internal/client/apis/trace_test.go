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

package apis

import (
	"internal/client/clienttest"
	"testing"
)

func TestListTracceSession(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := CreateProxy(proxyName, ""); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeployProxy(proxyName, 1, false, false, false, ""); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListTraceSession(proxyName, 1); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCreateTraceSession(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := CreateTraceSession(proxyName, 1, nil); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCleanupTrace(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := UndeployProxy(proxyName, 1, false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteProxy(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

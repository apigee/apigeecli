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

import "testing"

const (
	kvmName = "test"
)

func TestCreateProxyKVM(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := CreateProxy(proxyName, ""); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := CreateProxyKVM(proxyName, kvmName, true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestListProxyKVM(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := ListProxyKVM(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDeleteProxyKVM(t *testing.T) {
	if err := setup(false); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteProxyKVM(proxyName, kvmName); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := DeleteProxy(proxyName); err != nil {
		t.Fatalf("%v", err)
	}
}

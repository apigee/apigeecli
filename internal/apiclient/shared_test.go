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

package apiclient

import (
	"testing"

	"internal/clilog"
)

func TestInit(t *testing.T) {
	clilog.Init(true, true, true)
	clilog.Debug.Println("Printing Debug")
	clilog.Info.Println("Printing Information")
	clilog.Warning.Println("Printing Warning")
	clilog.Error.Println("Printing Error")
}

func TestHttpGet(t *testing.T) {
	clilog.Init(true, true, true)
	_, err := HttpClient("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpPost(t *testing.T) {
	payload := "test"
	_, err := HttpClient("https://httpbin.org/post", payload)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpDelete(t *testing.T) {
	_, err := HttpClient("https://httpbin.org/delete", "", "DELETE")
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpInvalidParam(t *testing.T) {
	_, err := HttpClient("https://httpbin.org/delete", "", "DELTE")
	if err == nil {
		t.Fatal(err)
	}
}

func TestHttpInvalidNumberOfParams(t *testing.T) {
	_, err := HttpClient("https://httpbin.org/delete", "", "DELETE", "SOMETHING ELSE")
	if err == nil {
		t.Fatal(err)
	}
}

func TestDownloadResource(t *testing.T) {
	// download 1000 bytes
	clilog.Init(true, true, true)
	err := DownloadResource("https://httpbin.org/stream-bytes/1000", "test", ".zip", false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteJSONArrayToFile(t *testing.T) {
	clilog.Init(true, true, true)
	entityPayloadList := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	err := WriteByteArrayToFile("test.json", false, entityPayloadList)
	if err != nil {
		t.Fatal(err)
	}
}

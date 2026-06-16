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
	"bytes"
	"internal/clilog"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newErrorResponse builds a minimal *http.Response carrying an error status
// code, mirroring what Apigee returns for e.g. a "not found" probe.
func newErrorResponse(statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"code":404,"message":"not found"}}`)),
		Header:     http.Header{},
	}
}

// TestHandleResponseSuppressesErrorWhenPrintDisabled reproduces
// https://github.com/apigee/apigeecli/issues/203: callers (like the
// products upsert pre-check GET) disable ClientPrintHttpResponse while
// probing whether an entity exists. The expected 404 from that probe must
// stay silent, the same way a successful probe response would.
func TestHandleResponseSuppressesErrorWhenPrintDisabled(t *testing.T) {
	clilog.Debug = log.New(io.Discard, "", 0)
	var buf bytes.Buffer
	clilog.HttpError = log.New(&buf, "", 0)

	ClientPrintHttpResponse.Set(false)
	defer ClientPrintHttpResponse.Set(true)

	_, err := handleResponse(newErrorResponse(404))
	if err == nil {
		t.Fatal("expected an error for a 404 response")
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no error output to be printed while ClientPrintHttpResponse is disabled, got: %q", buf.String())
	}
}

// TestHandleResponsePrintsErrorWhenPrintEnabled ensures genuine errors are
// still surfaced when the print toggle is left at its default (enabled).
func TestHandleResponsePrintsErrorWhenPrintEnabled(t *testing.T) {
	clilog.Debug = log.New(io.Discard, "", 0)
	var buf bytes.Buffer
	clilog.HttpError = log.New(&buf, "", 0)

	ClientPrintHttpResponse.Set(true)

	_, err := handleResponse(newErrorResponse(404))
	if err == nil {
		t.Fatal("expected an error for a 404 response")
	}
	if buf.Len() == 0 {
		t.Fatal("expected the error response to be printed when ClientPrintHttpResponse is enabled")
	}
}

// TestUpsertExistenceProbeStaysSilent exercises the real upsert pre-check
// flow end-to-end over HTTP: a GET against a product that doesn't exist
// (404), followed by the actual create POST (201). It mirrors exactly what
// products.upsert() does, including the ClientPrintHttpResponse toggling,
// to confirm the probe's expected 404 never reaches output while the
// subsequent create response still does.
func TestUpsertExistenceProbeStaysSilent(t *testing.T) {
	NewApigeeClient(ApigeeClientOptions{PrintOutput: true})
	SetApigeeToken("fake-test-token") // bypass the real auth flow

	mux := http.NewServeMux()
	mux.HandleFunc("/apiproducts/TestImport", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method %s on existence-check URL", r.Method)
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"ApiProduct with name TestImport does not exist"}}`))
	})
	mux.HandleFunc("/apiproducts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method %s on create URL", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"name":"TestImport"}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	var buf bytes.Buffer
	clilog.HttpError = log.New(&buf, "", 0)
	clilog.HttpResponse = log.New(&buf, "", 0)
	clilog.Debug = log.New(io.Discard, "", 0)

	// --- mirrors products.upsert() for Action == UPSERT ---
	ClientPrintHttpResponse.Set(false)
	_, err := HttpClient(server.URL + "/apiproducts/TestImport") // existence-check GET -> 404
	productMissing := err != nil
	ClientPrintHttpResponse.Set(GetCmdPrintHttpResponseSetting())

	if !productMissing {
		t.Fatal("expected the probe to report the product missing (err != nil)")
	}
	if buf.Len() != 0 {
		t.Fatalf("probe's 404 leaked to output even though ClientPrintHttpResponse was disabled: %q", buf.String())
	}

	respBody, err := HttpClient(server.URL+"/apiproducts", `{"name":"TestImport"}`) // real create POST -> 201
	if err != nil {
		t.Fatalf("unexpected error creating product: %v", err)
	}
	if len(respBody) == 0 {
		t.Fatal("expected a response body from the create call")
	}
	if buf.Len() == 0 {
		t.Fatal("expected the successful create response to be printed once the print toggle is restored")
	}
}

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

package keyaliases

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/apigee/apigeecli/apiclient"
)

func CreateSelfSigned(keystoreName string, name string, ignoreExpiry bool, ignoreNewLine bool, payload string) (respBody []byte, err error) {

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases")

	q := u.Query()
	q.Set("format", "selfsignedcert")
	q.Set("alias", name)

	if ignoreNewLine {
		q.Set("ignoreNewlineValidation", "true")
	}
	if ignoreExpiry {
		q.Set("ignoreExpiryValidation", "true")
	}
	u.RawQuery = q.Encode()

	var jsonPayload map[string]interface{}
	err = json.Unmarshal([]byte(payload), &jsonPayload)

	if err != nil {
		return respBody, err
	}

	return apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
}

func CreatePfx(keystoreName string, name string, ignoreExpiry bool, ignoreNewLine bool, pfxFile string, password string) (respBpdy []byte, err error) {

	if pfxFile == "" {
		return nil, fmt.Errorf("pfxFile cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	formParams := map[string]string{
		"file": pfxFile,
	}

	return create(keystoreName, name, "pkcs12", password, ignoreExpiry, ignoreNewLine, formParams)
}

func CreateKeyCert(keystoreName string, name string, ignoreExpiry bool, ignoreNewLine bool,
	certFile string, keyFile string, password string) (respBpdy []byte, err error) {

	if certFile == "" {
		return nil, fmt.Errorf("certFile cannot be empty")
	}

	formParams := map[string]string{
		"certFile": certFile,
	}
	if keyFile != "" {
		formParams["keyFile"] = keyFile
	}

	return create(keystoreName, name, "keycertfile", password, ignoreExpiry, ignoreNewLine, formParams)
}

func create(keystoreName string, name string, format string, password string, ignoreExpiry bool, ignoreNewLine bool, formParams map[string]string) (respBody []byte, err error) {

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases")

	q := u.Query()
	q.Set("format", format)
	q.Set("alias", name)

	if ignoreNewLine {
		q.Set("ignoreNewlineValidation", "true")
	}
	if ignoreExpiry {
		q.Set("ignoreExpiryValidation", "true")
	}
	if password != "" {
		q.Set("password", password)
	}
	u.RawQuery = q.Encode()

	return apiclient.PostHttpOctet(true, false, u.String(), formParams)
}

// CreateCSR
func CreateCSR(keystoreName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases", name, "csr")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return
}

// GetCert
func GetCert(keystoreName string, name string) (err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases", name, "certificate")
	err = apiclient.DownloadResource(u.String(), name+".crt", "")
	return err
}

// Get
func Get(keystoreName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

// Delete
func Delete(keystoreName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores",
		keystoreName, "aliases", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

// List
func List(keystoreName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

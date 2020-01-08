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

	"github.com/srinandan/apigeecli/apiclient"
)

func Create(keystoreName string, name string, format string, password string, ignoreExpiry bool,
	ignoreNewLine bool, payload string) (respBody []byte, err error) {
	if !validate(format) {
		return respBody, fmt.Errorf("certificate format must be one of keycertfile, pkcs12 or selfsignedcert")
	}
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

	switch format {
	case "keycertfile":
		if password != "" {
			q.Set("password", password)
		}
		u.RawQuery = q.Encode()
		respBody, err = apiclient.PostHttpOctet(true, u.String(), name+".pem")
	case "pkcs12":
		if password != "" {
			q.Set("password", password)
		}
		u.RawQuery = q.Encode()
		respBody, err = apiclient.PostHttpOctet(true, u.String(), name+".pfx")
	case "selfsignedcert":
		var jsonPayload map[string]interface{}
		err = json.Unmarshal([]byte(payload), &jsonPayload)

		if err != nil {
			return respBody, err
		}
		u.RawQuery = q.Encode()
		respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	}

	return respBody, err
}

func CreateCSR(keystoreName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases", name, "csr")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return
}

func GetCert(keystoreName string, name string) (err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases", name, "certificate")
	err = apiclient.DownloadResource(u.String(), name+".crt", "")
	return err
}

func Get(keystoreName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

func Delete(keystoreName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores",
		keystoreName, "aliases", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

func List(keystoreName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

func validate(format string) bool {
	var certFormats = [3]string{"keycertfile", "pkcs12", "selfsignedcert"}
	for _, frmt := range certFormats {
		if format == frmt {
			return true
		}
	}
	return false
}

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

	"internal/apiclient"

	"github.com/apigee/apigeecli/cmd/utils"
)

type certificate struct {
	Subject                    subject   `json:"subject" binding:"required"`
	KeySize                    *string   `json:"keySize" binding:"required"`
	SigAlg                     *string   `json:"sigAlg,omitempty"`
	SubjectAlternativeDNSNames *[]string `json:"subjectAlternativeDNSNames,omitempty"`
	CertValidityInDays         *string   `json:"certValidityInDays,omitempty"`
}

type subject struct {
	CountryCode *string `json:"countryCode,omitempty"`
	State       *string `json:"state,omitempty"`
	Locality    *string `json:"locality,omitempty"`
	Org         *string `json:"org,omitempty"`
	OrgUnit     *string `json:"orgUnit,omitempty"`
	CommonName  *string `json:"commonName" binding:"required"`
	Email       *string `json:"email,omitempty"`
}

func CreateOrUpdateSelfSigned(keystoreName string, name string, update bool, ignoreExpiry bool, ignoreNewLine bool, selfsignedFile string) (respBody []byte, err error) {
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

	payload, err := utils.ReadFile(selfsignedFile)
	if err != nil {
		return nil, err
	}

	cert := certificate{}
	err = json.Unmarshal([]byte(payload), &cert)
	if err != nil {
		return respBody, err
	}

	if cert.SigAlg == nil || *cert.SigAlg == "" {
		cert.SigAlg = new(string)
		*cert.SigAlg = "SHA256withRSA"
	}

	if cert.KeySize == nil || *cert.KeySize == "" {
		cert.KeySize = new(string)
		*cert.KeySize = "2048"
	}

	if cert.CertValidityInDays == nil || *cert.CertValidityInDays == "" {
		cert.CertValidityInDays = new(string)
		*cert.CertValidityInDays = "365"
	}

	if cert.Subject.CommonName == nil || *cert.Subject.CommonName == "" {
		return nil, fmt.Errorf("commonName is a mandatory parameter")
	}

	payload, err = json.Marshal(cert)
	if err != nil {
		return nil, err
	}

	if update {
		return apiclient.HttpClient(u.String(), string(payload), "PUT")
	}
	return apiclient.HttpClient(u.String(), string(payload))
}

func CreateOrUpdatePfx(keystoreName string, name string, update bool, ignoreExpiry bool, ignoreNewLine bool, pfxFile string, password string) (respBpdy []byte, err error) {
	if pfxFile == "" {
		return nil, fmt.Errorf("pfxFile cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	formParams := map[string]string{
		"file": pfxFile,
	}

	return createOrUpdate(keystoreName, name, "pkcs12", password, update, ignoreExpiry, ignoreNewLine, formParams)
}

func CreateOrUpdateKeyCert(keystoreName string, name string, update bool, ignoreExpiry bool, ignoreNewLine bool,
	certFile string, keyFile string, password string,
) (respBpdy []byte, err error) {
	if certFile == "" {
		return nil, fmt.Errorf("certFile cannot be empty")
	}

	formParams := map[string]string{
		"certFile": certFile,
	}
	if keyFile != "" {
		formParams["keyFile"] = keyFile
	}

	return createOrUpdate(keystoreName, name, "keycertfile", password, update, ignoreExpiry, ignoreNewLine, formParams)
}

func createOrUpdate(keystoreName string, name string, format string, password string, update bool, ignoreExpiry bool, ignoreNewLine bool, formParams map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	q := u.Query()
	q.Set("format", format)

	if update {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
			"keystores", keystoreName, "aliases", name)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
			"keystores", keystoreName, "aliases")
		q.Set("alias", name)
	}

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

	return apiclient.PostHttpOctet(update, u.String(), formParams)
}

// CreateCSR
func CreateCSR(keystoreName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases", name, "csr")
	respBody, err = apiclient.HttpClient(u.String())
	return
}

// GetCert
func GetCert(keystoreName string, name string) (err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases", name, "certificate")
	err = apiclient.DownloadResource(u.String(), name+".crt", "", true)
	return err
}

// Get
func Get(keystoreName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(keystoreName string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores",
		keystoreName, "aliases", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List(keystoreName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(),
		"keystores", keystoreName, "aliases")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

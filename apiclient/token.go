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
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/user"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwt"
	"github.com/srinandan/apigeecli/clilog"
)

const accessTokenFile = ".access_token"

type serviceAccount struct {
	Type                string `json:"type,omitempty"`
	ProjectID           string `json:"project_id,omitempty"`
	PrivateKeyID        string `json:"private_key_id,omitempty"`
	PrivateKey          string `json:"private_key,omitempty"`
	ClientEmail         string `json:"client_email,omitempty"`
	ClientID            string `json:"client_id,omitempty"`
	AuthURI             string `json:"auth_uri,omitempty"`
	TokenURI            string `json:"token_uri,omitempty"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientCertURL       string `json:"client_x509_cert_url,omitempty"`
}

var account = serviceAccount{}

func getPrivateKey(privateKey string) (interface{}, error) {
	pemPrivateKey := fmt.Sprintf("%v", privateKey)
	block, _ := pem.Decode([]byte(pemPrivateKey))
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		clilog.Error.Println("error parsing Private Key: ", err)
		return nil, err
	}
	return privKey, nil
}

func generateJWT(privateKey string) (string, error) {
	const aud = "https://www.googleapis.com/oauth2/v4/token"
	const scope = "https://www.googleapis.com/auth/cloud-platform"

	privKey, err := getPrivateKey(privateKey)

	if err != nil {
		return "", err
	}

	now := time.Now()
	token := jwt.New()

	_ = token.Set(jwt.AudienceKey, aud)
	_ = token.Set(jwt.IssuerKey, getServiceAccountProperty("ClientEmail"))
	_ = token.Set("scope", scope)
	_ = token.Set(jwt.IssuedAtKey, now.Unix())
	_ = token.Set(jwt.ExpirationKey, now.Unix())

	payload, err := token.Sign(jwa.RS256, privKey)
	if err != nil {
		clilog.Error.Println("error parsing Private Key: ", err)
		return "", err
	}
	clilog.Info.Println("jwt token : ", string(payload))
	return string(payload), nil
}

//generateAccessToken generates a Google OAuth access token from a service account
func generateAccessToken(privateKey string) (string, error) {

	const tokenEndpoint = "https://www.googleapis.com/oauth2/v4/token"
	const grantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"

	//oAuthAccessToken is a structure to hold OAuth response
	type oAuthAccessToken struct {
		AccessToken string `json:"access_token,omitempty"`
		ExpiresIn   int    `json:"expires_in,omitempty"`
		TokenType   string `json:"token_type,omitempty"`
	}

	token, err := generateJWT(privateKey)

	if err != nil {
		return "", nil
	}

	form := url.Values{}
	form.Add("grant_type", grantType)
	form.Add("assertion", token)

	client := &http.Client{}
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	resp, err := client.Do(req)

	if err != nil {
		clilog.Error.Println("failed to generate oauth token: ", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		clilog.Error.Println("error in response: ", string(bodyBytes))
		return "", errors.New("error in response")
	}
	decoder := json.NewDecoder(resp.Body)
	accessToken := oAuthAccessToken{}
	if err := decoder.Decode(&accessToken); err != nil {
		clilog.Error.Println("error in response: ", err)
		return "", errors.New("error in response")
	}
	clilog.Info.Println("access token : ", accessToken)
	SetApigeeToken(accessToken.AccessToken)
	_ = writeAccessToken()
	return accessToken.AccessToken, nil
}

func readServiceAccount(serviceAccountPath string) error {
	content, err := ioutil.ReadFile(serviceAccountPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &account)
	if err != nil {
		return err
	}
	return nil
}

func getServiceAccountProperty(key string) (value string) {
	r := reflect.ValueOf(&account)
	field := reflect.Indirect(r).FieldByName(key)
	return field.String()
}

func readAccessToken() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(path.Join(usr.HomeDir, accessTokenFile))
	if err != nil {
		clilog.Info.Println("Cached access token was not found")
		return err
	}
	clilog.Info.Println("Using cached access token: ", string(content))
	SetApigeeToken(string(content))
	return nil
}

func writeAccessToken() error {
	if IsSkipCache() {
		return nil
	}

	usr, err := user.Current()
	if err != nil {
		clilog.Warning.Println(err)
		return err
	}
	clilog.Info.Println("Cache access token: ", GetApigeeToken())
	//don't append access token
	return WriteByteArrayToFile(path.Join(usr.HomeDir, accessTokenFile), false, []byte(GetApigeeToken()))
}

func checkAccessToken() bool {
	if IsSkipCheck() {
		clilog.Info.Println("skipping token validity")
		return true
	}

	const tokenInfo = "https://www.googleapis.com/oauth2/v1/tokeninfo"
	u, _ := url.Parse(tokenInfo)
	q := u.Query()
	q.Set("access_token", GetApigeeToken())
	u.RawQuery = q.Encode()

	client := &http.Client{}

	clilog.Info.Println("Connecting to : ", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		clilog.Error.Println("error in client:", err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		clilog.Error.Println("error connecting to token endpoint: ", err)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clilog.Error.Println("token info error: ", err)
		return false
	} else if resp.StatusCode != 200 {
		clilog.Error.Println("token expired: ", string(body))
		return false
	}
	clilog.Info.Println("Response: ", string(body))
	clilog.Info.Println("Reusing the cached token: ", GetApigeeToken())
	return true
}

//SetAccessToken read from cache or if not found or expired will generate a new one
func SetAccessToken() error {
	if GetApigeeToken() == "" && GetServiceAccount() == "" {
		err := readAccessToken() //try to read from config
		if err != nil {
			return fmt.Errorf("either token or service account must be provided")
		}
		if checkAccessToken() { //check if the token is still valid
			return nil
		}
		return fmt.Errorf("token expired: request a new access token or pass the service account")
	}
	if GetServiceAccount() != "" {
		err := readServiceAccount(GetServiceAccount())
		if err != nil { // Handle errors reading the config file
			return fmt.Errorf("error reading config file: %s", err)
		}
		privateKey := getServiceAccountProperty("PrivateKey")
		if privateKey == "" {
			return fmt.Errorf("private key missing in the service account")
		}
		if getServiceAccountProperty("ClientEmail") == "" {
			return fmt.Errorf("client email missing in the service account")
		}
		_, err = generateAccessToken(privateKey)
		if err != nil {
			return fmt.Errorf("fatal error generating access token: %s", err)
		}
		return nil
	}
	//a token was passed, cache it
	if checkAccessToken() {
		_ = writeAccessToken()
		return nil
	}
	return fmt.Errorf("token expired: request a new access token or pass the service account")
}

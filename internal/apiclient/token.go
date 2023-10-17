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
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"internal/clilog"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/oauth2/google"
)

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

const tokenUri = "https://www.googleapis.com/oauth2/v4/token"

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
	const scope = "https://www.googleapis.com/auth/cloud-platform"

	privKey, err := getPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	now := time.Now()

	// Google OAuth takes aud as a string, not array
	// ref: https://github.com/lestrrat-go/jwx/releases/tag/v2.0.7
	jwt.Settings(jwt.WithFlattenAudience(true))
	token := jwt.New()
	token.Options().IsEnabled(jwt.FlattenAudience)

	_ = token.Set("aud", tokenUri)
	_ = token.Set(jwt.IssuerKey, getServiceAccountProperty("ClientEmail"))
	_ = token.Set("scope", scope)
	_ = token.Set(jwt.IssuedAtKey, now.Unix())
	_ = token.Set(jwt.ExpirationKey, now.Unix())

	payload, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privKey))
	if err != nil {
		clilog.Error.Println("error parsing Private Key: ", err)
		return "", err
	}
	clilog.Debug.Println("jwt token : ", string(payload))
	return string(payload), nil
}

// generateAccessToken generates a Google OAuth access token from a service account
func generateAccessToken(privateKey string) (string, error) {
	const grantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	var respBody []byte

	// oAuthAccessToken is a structure to hold OAuth response
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
	req, err := http.NewRequest("POST", tokenUri, strings.NewReader(form.Encode()))
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

	if resp != nil {
		defer resp.Body.Close()
	}

	if resp == nil {
		clilog.Error.Println("error in response: Response was null")
		return "", errors.New("error in response: Response was null")
	}

	respBody, err = io.ReadAll(resp.Body)
	clilog.Debug.Printf("Response: %s\n", string(respBody))

	if err != nil {
		clilog.Error.Println("error in response: ", err)
		return "", fmt.Errorf("error in response: %v", err)
	} else if resp.StatusCode > 399 {
		clilog.Error.Printf("status code %d, error in response: %s\n", resp.StatusCode, string(respBody))
		return "", fmt.Errorf("status code %d, error in response: %s\n", resp.StatusCode, string(respBody))
	}

	accessToken := oAuthAccessToken{}
	if err = json.Unmarshal(respBody, &accessToken); err != nil {
		return "", err
	}

	clilog.Debug.Println("access token : ", accessToken)

	SetApigeeToken(accessToken.AccessToken)
	_ = WriteToken(accessToken.AccessToken)
	return accessToken.AccessToken, nil
}

func readServiceAccount(serviceAccountPath string) error {
	content, err := os.ReadFile(serviceAccountPath)
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

func checkAccessToken() bool {
	if TokenCheckEnabled() {
		clilog.Debug.Println("skipping token validity")
		return true
	}

	const tokenInfo = "https://oauth2.googleapis.com/tokeninfo"
	u, _ := url.Parse(tokenInfo)
	q := u.Query()
	q.Set("access_token", GetApigeeToken())
	u.RawQuery = q.Encode()

	client := &http.Client{}

	clilog.Debug.Println("Connecting to : ", u.String())
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		clilog.Error.Println("token info error: ", err)
		return false
	} else if resp.StatusCode != 200 {
		clilog.Error.Println("token expired: ", string(body))
		return false
	}
	clilog.Debug.Println("Response: ", string(body))
	clilog.Debug.Println("Reusing the cached token: ", GetApigeeToken())
	return true
}

// SetAccessToken read from cache or if not found or expired will generate a new one
func SetAccessToken() error {
	if GetApigeeToken() == "" && GetServiceAccount() == "" {
		SetApigeeToken(GetToken()) // read from configuration
		if GetApigeeToken() == "" {
			return fmt.Errorf("either token or service account must be provided")
		}
		if checkAccessToken() { // check if the token is still valid
			return nil
		}
		return fmt.Errorf("token expired: request a new access token or pass the service account")
	}
	if GetApigeeToken() != "" {
		// a token was passed, cache it
		if checkAccessToken() {
			_ = WriteToken(GetApigeeToken())
			return nil
		}
	} else {
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
	return fmt.Errorf("token expired: request a new access token or pass the service account")
}

func getMetadata(metadata string) (respBpdy []byte, err error) {
	var req *http.Request

	metadataURL := fmt.Sprintf(
		"http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/%s", metadata)

	err = GetHttpClient()
	if err != nil {
		return nil, err
	}

	if DryRun() {
		return nil, nil
	}

	clilog.Debug.Println("Connecting to: ", metadataURL)

	req, err = http.NewRequest(http.MethodGet, metadataURL, nil)
	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return nil, err
	}

	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := ApigeeAPIClient.Do(req)
	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	if resp == nil {
		clilog.Error.Println("error in response: Response was null")
		return nil, fmt.Errorf("error in response: Response was null")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		clilog.Error.Println("error in response: ", err)
		return nil, err
	} else if resp.StatusCode > 399 {
		clilog.Debug.Printf("status code %d, error in response: %s\n", resp.StatusCode, string(respBody))
		clilog.HttpError.Println(string(respBody))
		return nil, errors.New(getErrorMessage(resp.StatusCode))
	}
	return respBody, err
}

// GetDefaultAccessToken
func GetDefaultAccessToken() (err error) {
	ctx := context.Background()
	tokenSource, err := google.DefaultTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return err
	}
	token, err := tokenSource.Token()
	if err != nil {
		return err
	}
	SetApigeeToken(token.AccessToken)
	return nil
}

// GetMetadataAccessToken
func GetMetadataAccessToken() (err error) {
	var tokenResponse map[string]interface{}

	respBody, err := getMetadata("token")
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, &tokenResponse)
	if err != nil {
		return err
	}

	SetApigeeToken(tokenResponse["access_token"].(string))

	ClientPrintHttpResponse.Set(false)
	defer ClientPrintHttpResponse.Set(GetCmdPrintHttpResponseSetting())

	respBody, _ = getMetadata("email")
	clilog.Debug.Println("service token email: ", string(respBody))

	respBody, _ = getMetadata("scopes")
	clilog.Debug.Println("scopes: ", string(respBody))

	return nil
}

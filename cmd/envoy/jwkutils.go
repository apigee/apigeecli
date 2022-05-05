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

package envoy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/apigee/apigeecli/clilog"
)

const keyFile = "remote-service.key"
const certFile = "remote-service.crt"
const kidFile = "remote-service.properties"
const use = "sig"
const kidFormat = "kid="

func readFile(name string) (data []byte, err error) {
	data, err = ioutil.ReadFile(name)
	return
}

func writeToFile(name string, data string) error {

	f, err := os.Create(name)
	if err != nil {
		log.Printf("failed to open file: %s\n", err)
		return err
	}

	_, err = f.WriteString(data)
	if err != nil {
		log.Printf("failed to open file: %s\n", err)
		return err
	}

	f.Close()

	return nil
}

func getPrivateKey(privateKey string) (interface{}, error) {
	pemPrivateKey := fmt.Sprintf("%v", privateKey)
	block, _ := pem.Decode([]byte(pemPrivateKey))
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		clilog.Error.Println("error parsing Private Key: ", err)
		return nil, err
	}
	return privKey, nil
}

func GenerateToken(folder string, expiry int) (string, error) {

	var jwtKidFile, jwtKeyFile, privateKey, kid string
	var err error

	if jwtKidFile, jwtKeyFile, _, err = checkFiles(folder); err != nil {
		return "", err
	}

	const aud = "remote-service-client"
	const iss = "apigee-remote-service-envoy"
	const tokenType = "JWT"

	if privateKey, err = getFileContents(jwtKeyFile); err != nil {
		return "", err
	}

	if kid, err = getFileContents(jwtKidFile); err != nil {
		return "", err
	}

	privKey, err := getPrivateKey(privateKey)

	if err != nil {
		return "", err
	}

	now := time.Now()
	hdr := jws.NewHeaders()
	if err = hdr.Set(jws.AlgorithmKey, jwa.RS256); err != nil {
		return "", err
	}
	if err = hdr.Set(jws.TypeKey, tokenType); err != nil {
		return "", err
	}

	if err = hdr.Set(jws.KeyIDKey, getKid(kid)); err != nil {
		return "", err
	}

	token := jwt.New()

	if err = token.Set(jwt.AudienceKey, aud); err != nil {
		return "", err
	}
	if err = token.Set(jwt.IssuerKey, iss); err != nil {
		return "", err
	}
	if err = token.Set(jwt.IssuedAtKey, now.Unix()); err != nil {
		return "", err
	}
	if err = token.Set(jwt.ExpirationKey, now.Unix()+(int64(expiry*60))); err != nil {
		return "", err
	}

	buf, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	payload, err := jws.Sign(buf, jwa.RS256, privKey, jws.WithHeaders(hdr))
	if err != nil {
		clilog.Error.Println("error parsing Private Key: ", err)
		return "", err
	}
	clilog.Info.Println("jwt token : ", string(payload))
	return string(payload), nil
}

func Generatekeys(kid string, folder string) (err error) {

	var jwtCertFile, jwtKeyFile string

	if folder != "" {
		jwtKeyFile = path.Join(folder, keyFile)
		jwtCertFile = path.Join(folder, certFile)
	} else {
		jwtKeyFile = keyFile
		jwtCertFile = certFile
	}

	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	)

	if err = writeToFile(jwtKeyFile, string(pemdata)); err != nil {
		return err
	}

	key, err := jwk.New(&privkey.PublicKey)
	if err != nil {
		return err
	}
	if err = key.Set(jwk.KeyUsageKey, use); err != nil {
		return err
	}
	if err = key.Set(jwk.KeyIDKey, kid); err != nil {
		return err
	}

	if fileExists(jwtCertFile) {
		//append cert
		return nil
	} else { //create new cert file
		jsonbuf, err := json.MarshalIndent(key, "", "  ")
		if err != nil {
			return err
		}

		set, err := jwk.ParseBytes(jsonbuf)
		if err != nil {
			return err
		}
		jsonbuf, err = json.MarshalIndent(set, "", "  ")
		if err != nil {
			return err
		}

		return writeToFile(jwtCertFile, string(jsonbuf))
	}
}

func Generatekid(kid string, folder string) (err error) {
	var jwtKidFile string
	if folder != "" {
		jwtKidFile = path.Join(folder, kidFile)
	} else {
		jwtKidFile = kidFile
	}
	data := kidFormat + kid
	return writeToFile(jwtKidFile, data)
}

func AddKey(kid string, folder string) (err error) {

	var jwtCertFile string
	if jwtCertFile = path.Join(folder, certFile); !fileExists(jwtCertFile) {
		return fmt.Errorf("remote-service.crt not found in %s", folder)
	}

	data, err := readFile(jwtCertFile)
	if err != nil {
		return err
	}

	set, err := jwk.ParseBytes(data)
	if err != nil {
		return err
	}

	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	)

	if err = writeToFile(keyFile, string(pemdata)); err != nil {
		return err
	}

	newKey, err := jwk.New(&privkey.PublicKey)
	if err != nil {
		return err
	}
	if err = newKey.Set(jwk.KeyUsageKey, use); err != nil {
		return err
	}
	if err = newKey.Set(jwk.KeyIDKey, kid); err != nil {
		return err
	}

	set.Keys = append(set.Keys, newKey)

	jsonbuf, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		return err
	}

	return writeToFile(certFile, string(jsonbuf))
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func getFileContents(filename string) (content string, err error) {
	var contentBytes []byte
	if contentBytes, err = ioutil.ReadFile(filename); err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(contentBytes), "\n"), nil
}

func getKid(kid string) string {
	return strings.ReplaceAll(kid, kidFormat, "")
}

func checkFiles(folder string) (jwtKidFile string, jwtKeyFile string, jwtCertFile string, err error) {

	if jwtKidFile = path.Join(folder, kidFile); !fileExists(jwtKidFile) {
		return "", "", "", fmt.Errorf("remote-service.properties not found in %s", folder)
	}

	if jwtKeyFile = path.Join(folder, keyFile); !fileExists(jwtKeyFile) {
		return "", "", "", fmt.Errorf("remote-service.key not found in %s", folder)
	}

	if jwtCertFile = path.Join(folder, certFile); !fileExists(jwtCertFile) {
		return "", "", "", fmt.Errorf("remote-service.crt not found in %s", folder)
	}

	return jwtKidFile, jwtKeyFile, jwtCertFile, nil
}

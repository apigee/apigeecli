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
	"log"
	"os"

	"github.com/lestrrat-go/jwx/jwk"
)

const keyFile = "remote-service.key"
const certFile = "remote-service.crt"
const kidFile = "remote-service.properties"

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

func Generatekeys(kid string) (err error) {

	const use = "sig"

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

	jsonbuf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return err
	}

	jsonWebKeys := "{ \"keys\":[" + string(jsonbuf) + "]}"
	return writeToFile(certFile, jsonWebKeys)

}

func Generatekid(kid string) (err error) {
	data := "kid=" + kid
	return writeToFile(kidFile, data)
}

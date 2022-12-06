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
	"fmt"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/keyaliases"
	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// CreateCmd to create key aliases
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Alias from PEM or PKCS12 file",
	Long:  "Create a Key Alias from PEM or PKCS12 file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		if format == "pfx" && password == "" {
			return fmt.Errorf("password must be set for pfx files")
		}
		if keyFile != "" && !utils.FileExists(keyFile) {
			return fmt.Errorf("keyFile was not found")
		}
		if certFile != "" && !utils.FileExists(certFile) {
			return fmt.Errorf("certFile was not found")
		}
		if pfxFile != "" && !utils.FileExists(pfxFile) {
			return fmt.Errorf("pfxFile was not found")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		switch format {
		case "pem":
			_, err = keyaliases.CreateKeyCert(keystoreName, name, ignoreExpiry, ignoreNewLine, certFile, keyFile, password)
		case "pkcs12":
			_, err = keyaliases.CreatePfx(keystoreName, name, ignoreExpiry, ignoreNewLine, pfxFile, password)
		default:
			return fmt.Errorf("invalid format key alias for %s", format)
		}
		return
	},
}

var format, password, keyFile, certFile, pfxFile string
var ignoreNewLine, ignoreExpiry bool

func init() {

	CreateCmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")
	CreateCmd.Flags().StringVarP(&name, "alias", "s",
		"", "Name of the key alias")
	CreateCmd.Flags().StringVarP(&format, "format", "f",
		"", "Format of the certificate; pem or pkcs12 (file extn is .pfx)")
	CreateCmd.Flags().StringVarP(&password, "password", "p",
		"", "PKCS12 password")
	CreateCmd.Flags().BoolVarP(&ignoreExpiry, "exp", "x",
		false, "Ignore expiry validation")
	CreateCmd.Flags().BoolVarP(&ignoreNewLine, "nl", "w",
		false, "Ignore new line in cert chain")
	CreateCmd.Flags().StringVarP(&certFile, "certFilePath", "",
		"", "Path to the X509 certificate in PEM format")
	CreateCmd.Flags().StringVarP(&keyFile, "keyFilePath", "",
		"", "Path to the X509 key in PEM format")
	CreateCmd.Flags().StringVarP(&pfxFile, "pfxFilePath", "",
		"", "Path to the PFX file")

	_ = CreateCmd.MarkFlagRequired("alias")
	_ = CreateCmd.MarkFlagRequired("format")
	_ = CreateCmd.MarkFlagRequired("key")
}

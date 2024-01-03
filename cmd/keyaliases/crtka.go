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

	"internal/apiclient"

	"internal/client/keyaliases"

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// CreateCmd to create key aliases
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Alias from PEM, PKCS12 or generate self signed cert",
	Long:  "Create a Key Alias from PEM, PKCS12 or generate self signed cert",
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
		if selfFile != "" && !utils.FileExists(selfFile) {
			return fmt.Errorf("selfsigned JSON file was not found")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		switch format {
		case "selfsignedcert":
			_, err = keyaliases.CreateOrUpdateSelfSigned(keystoreName,
				name,
				false,
				ignoreExpiry,
				ignoreNewLine,
				selfFile)
		case "keycertfile", "pem":
			_, err = keyaliases.CreateOrUpdateKeyCert(keystoreName,
				name,
				false,
				ignoreExpiry,
				ignoreNewLine,
				certFile,
				keyFile,
				password)
		case "pkcs12":
			_, err = keyaliases.CreateOrUpdatePfx(keystoreName,
				name,
				false,
				ignoreExpiry,
				ignoreNewLine,
				pfxFile,
				password)
		default:
			return fmt.Errorf("invalid format key alias for %s", format)
		}
		return err
	},
}

var (
	format, password, keyFile, certFile, pfxFile, selfFile string
	ignoreNewLine, ignoreExpiry                            bool
)

func init() {
	CreateCmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")
	CreateCmd.Flags().StringVarP(&name, "alias", "s",
		"", "Name of the key alias")
	CreateCmd.Flags().StringVarP(&format, "format", "f",
		"", "Format of the certificate; selfsignedcert, keycertfile (a.k.a pem), or pkcs12 (file extn is .pfx)")
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
	CreateCmd.Flags().StringVarP(&selfFile, "selfsignedFilePath", "",
		"", "Path to a JSON file containing details for a self signed certificate")

	_ = CreateCmd.MarkFlagRequired("alias")
	_ = CreateCmd.MarkFlagRequired("format")
	_ = CreateCmd.MarkFlagRequired("key")
}

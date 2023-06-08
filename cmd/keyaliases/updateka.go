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

// UpdateCmd to create key aliases
var UpdateCmd = &cobra.Command{
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
		if selfFile != "" && !utils.FileExists(selfFile) {
			return fmt.Errorf("selfsigned JSON file was not found")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		switch format {
		case "selfsignedcert":
			_, err = keyaliases.CreateOrUpdateSelfSigned(keystoreName,
				name,
				true,
				ignoreExpiry,
				ignoreNewLine,
				selfFile)
		case "pem":
			_, err = keyaliases.CreateOrUpdateKeyCert(keystoreName,
				name,
				true,
				ignoreExpiry,
				ignoreNewLine,
				certFile,
				keyFile,
				password)
		case "pkcs12":
			_, err = keyaliases.CreateOrUpdatePfx(keystoreName,
				name,
				true,
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

func init() {
	UpdateCmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")
	UpdateCmd.Flags().StringVarP(&name, "alias", "s",
		"", "Name of the key alias")
	UpdateCmd.Flags().StringVarP(&format, "format", "f",
		"", "Format of the certificate; selfsignedcert, pem or pkcs12 (file extn is .pfx)")
	UpdateCmd.Flags().StringVarP(&password, "password", "p",
		"", "PKCS12 password")
	UpdateCmd.Flags().BoolVarP(&ignoreExpiry, "exp", "x",
		false, "Ignore expiry validation")
	UpdateCmd.Flags().BoolVarP(&ignoreNewLine, "nl", "w",
		false, "Ignore new line in cert chain")
	UpdateCmd.Flags().StringVarP(&certFile, "certFilePath", "",
		"", "Path to the X509 certificate in PEM format")
	UpdateCmd.Flags().StringVarP(&keyFile, "keyFilePath", "",
		"", "Path to the X509 key in PEM format")
	UpdateCmd.Flags().StringVarP(&pfxFile, "pfxFilePath", "",
		"", "Path to the PFX file")
	UpdateCmd.Flags().StringVarP(&selfFile, "selfsignedFilePath", "",
		"", "Path to a JSON file containing details for a self signed certificate")

	_ = UpdateCmd.MarkFlagRequired("alias")
	_ = UpdateCmd.MarkFlagRequired("format")
	_ = UpdateCmd.MarkFlagRequired("key")
}

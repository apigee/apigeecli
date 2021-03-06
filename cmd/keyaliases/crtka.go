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
	"os"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/keyaliases"
	"github.com/apigee/apigeecli/clilog"
	"github.com/spf13/cobra"
)

//CreateCmd to create key aliases
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Alias from PEM or PKCS12 file",
	Long:  "Create a Key Alias from PEM or PKCS12 file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var fileformat string

		switch format {
		case "pem":
			fileformat = "pem"
			format = "keycertfile"
		case "pkcs12":
			fileformat = "pfx"
			format = "pkcs12"
		default:
			return fmt.Errorf("invalid format key alias for %s", format)
		}
		if !fileExists(name + "." + fileformat) {
			return fmt.Errorf("file %s.%s not found. Do not specify file extension with key alias name", name, fileformat)
		}
		_, err = keyaliases.Create(keystoreName, name, format, password, ignoreExpiry, ignoreNewLine, "")
		return
	},
}

var format, password string
var ignoreNewLine, ignoreExpiry bool

func init() {

	CreateCmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")
	CreateCmd.Flags().StringVarP(&name, "alias", "s",
		"", "Name of the key alias file. File name of pem or pfx file. Do not add the extension")
	CreateCmd.Flags().StringVarP(&format, "format", "f",
		"", "Format of the certificate; pem or pkcs12 (file extn is .pfx)")
	CreateCmd.Flags().StringVarP(&password, "password", "p",
		"", "PKCS12 password")
	CreateCmd.Flags().BoolVarP(&ignoreExpiry, "exp", "x",
		false, "Ignore expiry validation")
	CreateCmd.Flags().BoolVarP(&ignoreNewLine, "nl", "w",
		false, "Ignore new line in cert chain")

	_ = CreateCmd.MarkFlagRequired("alias")
	_ = CreateCmd.MarkFlagRequired("format")
	_ = CreateCmd.MarkFlagRequired("key")
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		clilog.Info.Println(err)
		return false
	}
	return true
}

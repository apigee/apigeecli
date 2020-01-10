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
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/keyaliases"
)

//Cmd to create key aliases
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Alias from PEM, JAR or PKCS12 file",
	Long:  "Create a Key Alias from PEM, JAR or PKCS12 file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		apiclient.SetApigeeEnv(env)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
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
		"", "Name of the key alias")
	CreateCmd.Flags().StringVarP(&format, "format", "f",
		"", "Format of the certificate")
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

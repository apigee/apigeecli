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
var CrtslfCmd = &cobra.Command{
	Use:   "create-self-signed",
	Short: "Create a Key Alias from self-seigned cert",
	Long:  "Create a Key Alias by generating a self-signed cert",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		apiclient.SetApigeeEnv(env)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = keyaliases.Create(keystoreName, aliasName, "selfsignedcert", "", ignoreExpiry, ignoreNewLine, payload)
		return
	},
}

const certFormat = "selfsignedcert"

var payload string

func init() {
	CrtslfCmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")
	CrtslfCmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key alias")
	CrtslfCmd.Flags().StringVarP(&payload, "cert", "c",
		"", "Certificate in JSON format")
	CrtslfCmd.Flags().BoolVarP(&ignoreExpiry, "exp", "x",
		false, "Ignore expiry validation")
	CrtslfCmd.Flags().BoolVarP(&ignoreNewLine, "nl", "w",
		false, "Ignore new line in cert chain")

	_ = CrtslfCmd.MarkFlagRequired("alias")
	_ = CrtslfCmd.MarkFlagRequired("cert")
	_ = CrtslfCmd.MarkFlagRequired("key")
}

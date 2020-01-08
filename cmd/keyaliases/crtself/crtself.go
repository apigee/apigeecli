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

package crtself

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/client/keyaliases"
)

//Cmd to create key aliases
var Cmd = &cobra.Command{
	Use:   "create-self-signed",
	Short: "Create a Key Alias from self-seigned cert",
	Long:  "Create a Key Alias by generating a self-signed cert",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = keyaliases.Create(keystoreName, aliasName, "selfsignedcert", "", ignoreExpiry, ignoreNewLine, payload)
		return
	},
}

const certFormat = "selfsignedcert"

var keystoreName, aliasName, payload string
var ignoreNewLine, ignoreExpiry bool

func init() {
	Cmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")
	Cmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key alias")
	Cmd.Flags().StringVarP(&payload, "cert", "c",
		"", "Certificate in JSON format")
	Cmd.Flags().BoolVarP(&ignoreExpiry, "exp", "x",
		false, "Ignore expiry validation")
	Cmd.Flags().BoolVarP(&ignoreNewLine, "nl", "w",
		false, "Ignore new line in cert chain")

	_ = Cmd.MarkFlagRequired("alias")
	_ = Cmd.MarkFlagRequired("cert")
	_ = Cmd.MarkFlagRequired("key")
}

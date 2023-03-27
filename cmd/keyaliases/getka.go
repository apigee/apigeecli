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
	"internal/apiclient"

	"internal/client/keyaliases"

	"github.com/spf13/cobra"
)

// Cmd to get key aliases
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a Key Alias",
	Long:  "Get a Key Alias",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = keyaliases.Get(keystoreName, aliasName)
		return
	},
}

func init() {

	GetCmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")

	GetCmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key store")

	_ = GetCmd.MarkFlagRequired("alias")
	_ = GetCmd.MarkFlagRequired("key")
}

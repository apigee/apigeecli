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

package keystores

import (
	"internal/apiclient"

	"internal/client/keystores"

	"github.com/spf13/cobra"
)

// CreateCmd to create key stores
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Store",
	Long:  "Create a Key Store",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = keystores.Create(name)
		return
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the key store")

	_ = CreateCmd.MarkFlagRequired("name")
}

// Copyright 2023 Google LLC
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

package env

import (
	"internal/apiclient"
	"internal/client/env"

	"github.com/spf13/cobra"
)

// GetSecInCmd returns security incidents
var GetSecInCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns a security incidents by name",
	Long:  "Returns a security incidents by name",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = env.GetSecurityIncident(name)
		return
	},
}

func init() {
	GetSecInCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the security incident")
	_ = GetSecInCmd.MarkFlagRequired("name")
}

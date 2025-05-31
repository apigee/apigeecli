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

// ListSecActCmd to list catalog items
var ListSecActCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns security actions in the environment",
	Long:  "Returns security actions in the environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = env.ListSecurityActions(pageSize, pageToken, filter)
		return
	},
}

func init() {
	ListSecActCmd.Flags().IntVarP(&pageSize, "page-size", "",
		-1, "The maximum number of versions to return")
	ListSecActCmd.Flags().StringVarP(&pageToken, "page-token", "",
		"", "A page token, received from a previous call")
	ListSecActCmd.Flags().StringVarP(&filter, "filter", "",
		"", "Filter results")
}

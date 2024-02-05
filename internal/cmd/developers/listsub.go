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

package developers

import (
	"internal/apiclient"

	"internal/client/developers"

	"github.com/spf13/cobra"
)

// ListSubCmd to list developer subscriptions
var ListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of Developer subscriptions",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	Long: "Returns a list of Developer subscriptions",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = developers.ListSubscriptions(email)
		return
	},
}

func init() {
	ListSubCmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")

	_ = ListSubCmd.MarkFlagRequired("email")
}

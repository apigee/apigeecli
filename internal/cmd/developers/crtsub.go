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

// CreateSubCmd to create developer subscription
var CreateSubCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a subscription for a developer",
	Long:  "Create a subscription for a developer",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = developers.CreateSubscription(email, apiproduct)
		return
	},
}

var apiproduct string

func init() {
	CreateSubCmd.Flags().StringVarP(&email, "email", "e",
		"", "The developer's email")
	CreateSubCmd.Flags().StringVarP(&apiproduct, "apiproduct", "p",
		"", "The name of the API Product")

	_ = CreateSubCmd.MarkFlagRequired("email")
	_ = CreateSubCmd.MarkFlagRequired("apiproduct")
}

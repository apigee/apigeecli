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

package apps

import (
	"internal/apiclient"

	"internal/client/apps"

	"github.com/spf13/cobra"
)

// ManageKeyCmd to create developer keys
var ManageKeyCmd = &cobra.Command{
	Use:   "manage",
	Short: "Approve or revoke a developer app key",
	Long:  "Approve or revoke a developer app key",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apps.ManageKey(developerEmail, name, key, action)
		return
	},
}

var action string

func init() {
	ManageKeyCmd.Flags().StringVarP(&key, "key", "k",
		"", "Developer app consumer key")
	ManageKeyCmd.Flags().StringVarP(&action, "action", "x",
		"revoke", "Action to perform - revoke or approve")

	_ = ManageKeyCmd.MarkFlagRequired("key")
	_ = ManageKeyCmd.MarkFlagRequired("action")
}

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

// DelCmd to delete developer
var DelCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an App Developer from an organization",
	Long:  "Deletes an App Developer from an organization",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = developers.Delete(email)
		return
	},
}

func init() {
	DelCmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")

	_ = DelCmd.MarkFlagRequired("email")
}

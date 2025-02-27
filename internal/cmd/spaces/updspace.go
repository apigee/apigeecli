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

package spaces

import (
	"internal/apiclient"
	"internal/client/spaces"

	"github.com/spf13/cobra"
)

// UpdateCmd to update Apigee Space
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an Apigee Space",
	Long:  "Update an Apigee Space",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = spaces.Update(name, displayName)
		return err
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Space")
	UpdateCmd.Flags().StringVarP(&displayName, "displayName", "d",
		"", "Display Name for the Space")

	_ = UpdateCmd.MarkFlagRequired("name")
	_ = UpdateCmd.MarkFlagRequired("displayName")
}

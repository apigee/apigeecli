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

package securityprofiles

import (
	"internal/apiclient"
	"internal/client/securityprofiles"

	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// UpdateCmd to update a security profile
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing Security Profile",
	Long:  "Update an existing Security Profile",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		content, err := utils.ReadFile(securityActionFile)
		if err != nil {
			return err
		}
		_, err = securityprofiles.Update(name, content)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the security profile")
	UpdateCmd.Flags().StringVarP(&securityActionFile, "file", "f",
		"", "Path to a file containing Security Profile content")
	_ = UpdateCmd.MarkFlagRequired("name")
	_ = UpdateCmd.MarkFlagRequired("file")
}

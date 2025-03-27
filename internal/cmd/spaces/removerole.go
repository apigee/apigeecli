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
	"internal/clilog"

	"github.com/spf13/cobra"
)

// RemoveRoleCmd to a member from a role
var RemoveRoleCmd = &cobra.Command{
	Use:   "removerole",
	Short: "Remove a member or SA from a role for a space",
	Long:  "Remove a member or SA from a role for a space",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		err = spaces.RemoveIAM(space, memberName, role)
		if err != nil {
			return err
		}
		clilog.Info.Printf("Member \"%s\" removed access to role \"%s\" in space \"%s\"\n", memberName, role, space)
		return nil
	},
	Example: `Remove Space Editor role for user in a space: ` + GetExample(4),
}

func init() {
	RemoveRoleCmd.Flags().StringVarP(&space, "space", "",
		"", "Space name.")
	RemoveRoleCmd.Flags().StringVarP(&memberName, "name", "n",
		"", "Member name. Prefix the member role as role:name. Ex: user:foo@acme.com")
	RemoveRoleCmd.Flags().StringVarP(&role, "role", "",
		"", "IAM Role")

	_ = RemoveRoleCmd.MarkFlagRequired("space")
	_ = RemoveRoleCmd.MarkFlagRequired("name")
	_ = RemoveRoleCmd.MarkFlagRequired("role")
}

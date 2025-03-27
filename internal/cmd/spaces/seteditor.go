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

// SetEditorCmd to set role on env
var SetEditorCmd = &cobra.Command{
	Use:   "seteditor",
	Short: "Set Space Content Editor role for a member on a Space",
	Long:  "Set Space Content Editor role for a member a Space",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		err = spaces.SetIAM(space, memberName, "editor", memberType)
		if err != nil {
			return err
		}
		clilog.Info.Printf("Member \"%s\" granted access to \"roles/apigee.spaceContentEditor\" role in space \"%s\"\n", memberName, space)
		return nil
	},
	Example: `Set Space Editor role for user in a space: ` + GetExample(2),
}

func init() {
	SetEditorCmd.Flags().StringVarP(&space, "space", "",
		"", "Space name.")
	SetEditorCmd.Flags().StringVarP(&memberName, "name", "n",
		"", "Member Name, example Service Account Name")
	SetEditorCmd.Flags().StringVarP(&memberType, "member-type", "m",
		"serviceAccount", "memberType must be serviceAccount, user or group")

	_ = SetEditorCmd.MarkFlagRequired("space")
	_ = SetEditorCmd.MarkFlagRequired("name")
}


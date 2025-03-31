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

// SetViewerCmd to set role on env
var SetViewerCmd = &cobra.Command{
    Use:   "setviewer",
    Short: "Set Space Content Viewer role for a member on a Space",
    Long:  "Set Space Content Viewer role for a member on a Space",
    Args: func(cmd *cobra.Command, args []string) (err error) {
        apiclient.SetRegion(region)
        return apiclient.SetApigeeOrg(org)
    },
    RunE: func(cmd *cobra.Command, args []string) (err error) {
        cmd.SilenceUsage = true

        err = spaces.SetIAM(space, memberName, "viewer", memberType)
        if err != nil {
            return err
        }
        clilog.Info.Printf("Member \"%s\" granted access to \"roles/apigee.spaceContentViewer\" role in space \"%s\"\n", memberName, space)
        return nil
    },
    Example: `Set Space Viewer role for user in a space: ` + GetExample(3),
}

func init() {
    SetViewerCmd.Flags().StringVarP(&space, "space", "",
        "", "Space name.")
    SetViewerCmd.Flags().StringVarP(&memberName, "name", "n",
        "", "Member Name, example Service Account Name")
    SetViewerCmd.Flags().StringVarP(&memberType, "member-type", "m",
        "serviceAccount", "memberType must be serviceAccount, user or group")

    _ = SetViewerCmd.MarkFlagRequired("space")
    _ = SetViewerCmd.MarkFlagRequired("name")
}


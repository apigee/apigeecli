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

package env

import (
	"internal/apiclient"
	"internal/clilog"

	environments "internal/client/env"

	"github.com/spf13/cobra"
)

// SetAdminCmd to set role on env
var SetAdminCmd = &cobra.Command{
	Use:   "setadmin",
	Short: "Set Environment Admin role for a member on an environment",
	Long:  "Set Environment Admin role for a member an Environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		err = environments.SetIAM(memberName, "admin", memberType)
		if err != nil {
			return err
		}
		clilog.Info.Printf("Member %s granted access to environment admin role\n", memberName)
		return nil
	},
}

func init() {
	SetAdminCmd.Flags().StringVarP(&memberName, "name", "n",
		"", "Member Name, example Service Account Name")
	SetAdminCmd.Flags().StringVarP(&memberType, "member-type", "m",
		"serviceAccount", "memberType must be serviceAccount, user or group")

	_ = SetAdminCmd.MarkFlagRequired("name")
}

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
	"fmt"
	"regexp"

	"internal/apiclient"

	environments "internal/client/env"

	"github.com/spf13/cobra"
)

// SetCustCmd to manage custom roles for an env
var SetCustCmd = &cobra.Command{
	Use:   "setcustom",
	Short: "Set a custom role for a member on an environment",
	Long:  "Set a custom role for a member on an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		re := regexp.MustCompile(`projects\/([a-zA-Z0-9_-]+)\/roles\/([a-zA-Z0-9_-]+)`)
		result := re.FindString(role)
		if result == "" {
			return fmt.Errorf("custom role must be of the format projects/{project-id}/roles/{role-name}")
		}
		err = environments.SetIAM(memberName, role, memberType)
		if err != nil {
			return err
		}
		fmt.Printf("Member %s, granted access to %s\n", memberName, role)
		return nil
	},
}

func init() {

	SetCustCmd.Flags().StringVarP(&memberName, "name", "n",
		"", "Member Name, example Service Account Name")
	SetCustCmd.Flags().StringVarP(&role, "role", "r",
		"", "Custom IAM role in the format projects/{project-id}/roles/{role}")
	SetCustCmd.Flags().StringVarP(&memberType, "memberType", "m",
		"serviceAccount", "memberType must be serviceAccount, user or group")

	_ = SetCustCmd.MarkFlagRequired("name")
	_ = SetCustCmd.MarkFlagRequired("role")
}

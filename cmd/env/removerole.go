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

	"internal/apiclient"

	environments "github.com/apigee/apigeecli/client/env"
	"github.com/spf13/cobra"
)

// RemoveRoleCmd to a member from a role
var RemoveRoleCmd = &cobra.Command{
	Use:   "removerole",
	Short: "Remove a member or SA from a role for an environment",
	Long:  "Remove a member or SA from a role for an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = environments.RemoveIAM(memberName, role)
		if err != nil {
			return err
		}
		fmt.Printf("Member %s removed access to role %s\n", memberName, role)
		return nil
	},
}

func init() {

	RemoveRoleCmd.Flags().StringVarP(&memberName, "name", "n",
		"", "Member name. Prefix the member role as role:name. Ex: user:foo@acme.com")
	RemoveRoleCmd.Flags().StringVarP(&role, "role", "r",
		"", "IAM Role")

	_ = RemoveRoleCmd.MarkFlagRequired("name")
	_ = RemoveRoleCmd.MarkFlagRequired("role")
}

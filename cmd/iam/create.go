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

package iam

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
)

//Cmd to get org details
var CallCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new IAM Service Account with permissions for Apigee Runtime",
	Long:  "Create a new IAM Service Account with permissions for Apigee Runtime",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if !generateName && name == "" {
			return fmt.Errorf("provide a service account name or allow the tool to generate one")
		}
		if !ValidateRoleType(roleType) {
			return fmt.Errorf("The role type %s is not a valid type. Please use one of %s", roleType, roles)
		}
		apiclient.SetProjectID(projectID)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if generateName {
			name = GenerateName("apigee-" + roleType + "-")
		}
		return apiclient.CreateIAMServiceAccount(name, roleType)
	},
}

func init() {

	CallCmd.Flags().StringVarP(&projectID, "prj", "p",
		"", "GCP Project ID")
	CallCmd.Flags().StringVarP(&name, "name", "n",
		"", "Service Account Name")
	CallCmd.Flags().BoolVarP(&generateName, "gen", "g",
		false, "Generate account name")
	CallCmd.Flags().StringVarP(&roleType, "role", "r",
		"", "IAM Role Type")

	_ = CallCmd.MarkFlagRequired("prj")
	_ = CallCmd.MarkFlagRequired("role")
}

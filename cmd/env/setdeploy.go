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

	"github.com/apigee/apigeecli/apiclient"
	environments "github.com/apigee/apigeecli/client/env"
	"github.com/spf13/cobra"
)

//Cmd to manage tracing of apis
var SetDepCmd = &cobra.Command{
	Use:   "setdeploy",
	Short: "Set Apigee Deployer role for a member on an environment",
	Long:  "Set Apigee Deployer role for a member on an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = environments.SetIAM(memberName, "deploy", memberType)
		if err != nil {
			return err
		}
		fmt.Printf("Member %s granted access to Apigee Deployer role\n", memberName)
		return nil
	},
}

func init() {

	SetDepCmd.Flags().StringVarP(&memberName, "name", "n",
		"", "Member Name, example Service Account Name")
	SetDepCmd.Flags().StringVarP(&memberType, "memberType", "m",
		"serviceAccount", "memberType must be serviceAccount, user or group")

	_ = SetDepCmd.MarkFlagRequired("name")
}

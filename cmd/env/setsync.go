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

// Cmd to manage tracing of apis
var SetSyncCmd = &cobra.Command{
	Use:   "setsync",
	Short: "Set Synchronization Manager role for a member on an environment",
	Long:  "Set Synchronization Manager role for a member on an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = environments.SetIAM(memberName, "sync", memberType)
		if err != nil {
			return err
		}
		fmt.Printf("Member %s granted access to Apigee Synchronizer Manager role\n", memberName)
		return nil
	},
}

func init() {

	SetSyncCmd.Flags().StringVarP(&memberName, "name", "n",
		"", "Member Name, example Service Account Name")
	SetSyncCmd.Flags().StringVarP(&memberType, "memberType", "m",
		"serviceAccount", "memberType must be serviceAccount, user or group")

	_ = SetSyncCmd.MarkFlagRequired("name")
}

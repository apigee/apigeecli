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

package org

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/orgs"
)

//UpdateCmd to get org details
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update settings of an Apigee Org",
	Long:  "Update settings of an Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if runtimeType != "HYBRID" && runtimeType != "CLOUD" {
			return fmt.Errorf("runtime type must be CLOUD or HYBRID")
		}
		if runtimeType == "CLOUD" {
			if network == "" {
				return fmt.Errorf("authorized network must be supplied")
			}
			if databaseKey == "" {
				return fmt.Errorf("runtime database encryption key must be supplied")
			}
		}
		apiclient.SetProjectID(projectID)
		return apiclient.SetApigeeOrg(projectID)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = orgs.Update(description, "", region, network, runtimeType, databaseKey)
		return
	},
}

func init() {

	UpdateCmd.Flags().StringVarP(&region, "reg", "r",
		"", "Analytics region name")
	UpdateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Apigee org description")
	UpdateCmd.Flags().StringVarP(&network, "net", "n",
		"default", "Authorized network; if using a shared VPC format is projects/{host-project-id}/{location}/networks/{network-name}")
	UpdateCmd.Flags().StringVarP(&databaseKey, "key", "k",
		"", "Runtime Database Encryption Key")
	UpdateCmd.Flags().StringVarP(&runtimeType, "runtime-type", "",
		"HYBRID", "Runtime type: CLOUD or HYBRID")

	_ = UpdateCmd.MarkFlagRequired("prj")
	_ = UpdateCmd.MarkFlagRequired("reg")
	_ = UpdateCmd.MarkFlagRequired("runtime-type")
}

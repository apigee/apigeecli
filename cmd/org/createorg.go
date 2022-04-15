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

//Cmd to get org details
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Apigee Org",
	Long:  "Create a new Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if runtimeType != "HYBRID" && runtimeType != "CLOUD" {
			return fmt.Errorf("runtime type must be CLOUD or HYBRID")
		}

		if billingType != "SUBSCRIPTION" && billingType != "EVALUATION" {
			return fmt.Errorf("Billing type must be SUBSCRIPTION or EVALUATION")
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
		_, err = orgs.Create(region, network, runtimeType, databaseKey, billingType, disablePortal)
		return
	},
}

var region, projectID, network, runtimeType, description, databaseKey, billingType string
var disablePortal bool

func init() {

	CreateCmd.Flags().StringVarP(&region, "reg", "r",
		"", "Analytics region name")
	CreateCmd.Flags().StringVarP(&projectID, "prj", "p",
		"", "GCP Project ID")
	CreateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Apigee org description")
	CreateCmd.Flags().StringVarP(&network, "net", "n",
		"default", "Authorized network; if using a shared VPC format is projects/{host-project-id}/{location}/networks/{network-name}")
	CreateCmd.Flags().StringVarP(&databaseKey, "key", "k",
		"", "Runtime Database Encryption Key")
	CreateCmd.Flags().StringVarP(&runtimeType, "runtime-type", "",
		"HYBRID", "Runtime type: CLOUD or HYBRID")
	CreateCmd.Flags().StringVarP(&runtimeType, "billing-type", "",
		"", "Billing type: SUBSCRIPTION or EVALUATION")
	CreateCmd.Flags().BoolVarP(&disablePortal, "disable-portal", "",
		false, "Disable creation of Developer Portals")

	_ = CreateCmd.MarkFlagRequired("prj")
	_ = CreateCmd.MarkFlagRequired("reg")
	_ = CreateCmd.MarkFlagRequired("runtime-type")
}

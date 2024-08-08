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
	"encoding/json"
	"fmt"
	"path/filepath"

	"internal/apiclient"

	"internal/client/operations"
	"internal/client/orgs"

	"github.com/spf13/cobra"
)

const (
	HYBRID       = "HYBRID"
	CLOUD        = "CLOUD"
	SUBSCRIPTION = "SUBSCRIPTION"
	EVALUATION   = "EVALUATION"
)

// CreateCmd to get org details
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Apigee Org",
	Long:  "Create a new Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if runtimeType != HYBRID && runtimeType != CLOUD {
			return fmt.Errorf("runtime type must be %s or %s", CLOUD, HYBRID)
		}

		if billingType != SUBSCRIPTION && billingType != EVALUATION {
			return fmt.Errorf("billing type must be %s or %s", SUBSCRIPTION, EVALUATION)
		}
		apiclient.SetProjectID(projectID)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(projectID)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		respBody, err := orgs.Create(description, analyticsRegion, authorizedNetwork,
			disableVpcPeering, runtimeType, billingType, runtimeDatabaseEncryptionKeyName,
			portalDisabled, apiConsumerDataEncryptionKeyName, controlPlaneEncryptionKeyName,
			apiConsumerDataLocation)
		if err != nil {
			return
		}
		if wait {
			respMap := make(map[string]interface{})
			if err = json.Unmarshal(respBody, &respMap); err != nil {
				return err
			}
			err = operations.WaitForOperation(filepath.Base(respMap["name"].(string)))
		}
		return
	},
}

var (
	analyticsRegion, projectID, authorizedNetwork, runtimeType                               string
	description, runtimeDatabaseEncryptionKeyName, billingType                               string
	apiConsumerDataEncryptionKeyName, controlPlaneEncryptionKeyName, apiConsumerDataLocation string
	disableVpcPeering, portalDisabled, wait                                                  bool
)

func init() {
	CreateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Apigee org description")
	CreateCmd.Flags().StringVarP(&analyticsRegion, "reg", "",
		"", "Analytics region name")
	CreateCmd.Flags().BoolVarP(&disableVpcPeering, "disable-vpc-peering", "",
		true, "Disable VPC Peering; default true")
	CreateCmd.Flags().StringVarP(&runtimeType, "runtime-type", "",
		HYBRID, "Runtime type: CLOUD or HYBRID")
	CreateCmd.Flags().StringVarP(&runtimeType, "billing-type", "",
		"", "Billing type: SUBSCRIPTION or EVALUATION")
	CreateCmd.Flags().StringVarP(&runtimeDatabaseEncryptionKeyName, "runtime-key", "",
		"", "Runtime Database Encryption Key")
	CreateCmd.Flags().StringVarP(&apiConsumerDataEncryptionKeyName, "consumer-key", "",
		"", "API Consumer Data Encryption Key")
	CreateCmd.Flags().StringVarP(&controlPlaneEncryptionKeyName, "controlplane-key", "",
		"", "Apigee Controlplane Encryption Key")
	CreateCmd.Flags().StringVarP(&apiConsumerDataLocation, "key", "k",
		"", "API Consumer data location")

	CreateCmd.Flags().StringVarP(&projectID, "prj", "p",
		"", "GCP Project ID")

	CreateCmd.Flags().StringVarP(&authorizedNetwork, "net", "n",
		"default", "Authorized network; if using a shared VPC format is "+
			"projects/{host-project-id}/{location}/networks/{network-name}")

	CreateCmd.Flags().BoolVarP(&portalDisabled, "disable-portal", "",
		false, "Disable creation of Developer Portals; default false")

	CreateCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the create to finish, with success or error")

	_ = CreateCmd.MarkFlagRequired("prj")
	_ = CreateCmd.MarkFlagRequired("reg")
	_ = CreateCmd.MarkFlagRequired("runtime-type")
}

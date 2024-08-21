// Copyright 2024 Google LLC
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

package dependencies

import (
	"fmt"
	"internal/apiclient"
	"internal/client/hub"
	"regexp"

	"github.com/spf13/cobra"
)

// CrtCmd
var CrtCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Dependency in API Hub",
	Long:  "Create a new Dependency in API Hub",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		extApiPattern := `^projects/[^/]+/locations/[^/]+/apis/[^/]+/externalApis/[^/]+$`
		opPattern := `^projects/[^/]+/locations/[^/]+/apis/[^/]+/versions/[^/]+/operations/[^/]+$`

		if consumerExternalApiResouceName == "" && consumerOperationResourceName == "" {
			return fmt.Errorf("at least one of consumer-ext-api-res-name or consumer-op-res-name must be set")
		}
		if consumerExternalApiResouceName != "" && consumerOperationResourceName != "" {
			return fmt.Errorf("consumer-ext-api-res-name and consumer-op-res-name cannot be used together")
		}
		if supplierExternalApiResourceName == "" && supplierOperationResourceName == "" {
			return fmt.Errorf("at least one of supplier-ext-api-res-name and supplier-op-res-name must be set")
		}
		if supplierExternalApiResourceName != "" && supplierOperationResourceName != "" {
			return fmt.Errorf("supplier-ext-api-res-name and supplier-op-res-name cannot be used together")
		}

		if ok, err := regexp.MatchString(opPattern, consumerOperationResourceName); err != nil || !ok {
			return fmt.Errorf("invalid pattern for consumer-op-res-name: %v", err)
		}
		if ok, err := regexp.MatchString(extApiPattern, consumerExternalApiResouceName); err != nil || !ok {
			return fmt.Errorf("invalid pattern for consumer-ext-api-res-name: %v", err)
		}

		if ok, err := regexp.MatchString(opPattern, supplierOperationResourceName); err != nil || !ok {
			return fmt.Errorf("invalid pattern for supplier-op-res-name: %v", err)
		}
		if ok, err := regexp.MatchString(extApiPattern, supplierExternalApiResourceName); err != nil || !ok {
			return fmt.Errorf("invalid pattern for supplier-ext-api-res-name: %v", err)
		}

		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = hub.CreateDependency(dependencyID, description, consumerDisplayName,
			consumerOperationResourceName, consumerExternalApiResouceName, supplierDisplayName,
			supplierOperationResourceName, supplierExternalApiResourceName)
		return
	},
}

var (
	dependencyID, description, consumerDisplayName, consumerOperationResourceName      string
	consumerExternalApiResouceName, supplierDisplayName, supplierOperationResourceName string
	supplierExternalApiResourceName                                                    string
)

func init() {
	CrtCmd.Flags().StringVarP(&dependencyID, "id", "i",
		"", "Dependency Display Name")
	CrtCmd.Flags().StringVarP(&description, "description", "",
		"", "Dependency Description")
	CrtCmd.Flags().StringVarP(&consumerDisplayName, "consumer-display-name", "",
		"", "Consumer Display Name")
	CrtCmd.Flags().StringVarP(&consumerOperationResourceName, "consumer-op-res-name", "",
		"", "Consumer Operation Resource Name")
	CrtCmd.Flags().StringVarP(&consumerExternalApiResouceName, "consumer-ext-api-res-name", "",
		"", "Consumer External API Resource Name")
	CrtCmd.Flags().StringVarP(&supplierDisplayName, "supplier-display-name", "",
		"", "Supplier Display Name")
	CrtCmd.Flags().StringVarP(&supplierOperationResourceName, "supplier-op-res-name", "",
		"", "Supplier Operation Resource Name")
	CrtCmd.Flags().StringVarP(&supplierExternalApiResourceName, "supplier-ext-api-res-name", "",
		"", "Supplier External API Resource Name")

	_ = CrtCmd.MarkFlagRequired("id")
}

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

package externalapis

import (
	"internal/apiclient"
	"internal/client/hub"

	"github.com/spf13/cobra"
)

// UpdateCmd
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an External API",
	Long:  "Update an External API",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = hub.UpdateExternalAPI(externalApiID, displayName,
			description, endpoints, paths, externalURI, attribute)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&externalApiID, "id", "i",
		"", "External API ID")
	UpdateCmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "Exernal API Display Name")
	UpdateCmd.Flags().StringVarP(&description, "description", "",
		"", "External API Description")
	UpdateCmd.Flags().StringVarP(&externalURI, "external-uri", "",
		"", "The uri of the externally hosted documentation")
	UpdateCmd.Flags().StringArrayVarP(&endpoints, "endpoints", "",
		[]string{}, " The endpoints at which this deployment resource is listening for API requests")
	UpdateCmd.Flags().StringVarP(&attribute, "attribute", "",
		"", "The name of the attribute for the external api")
	UpdateCmd.Flags().StringArrayVarP(&paths, "paths", "",
		[]string{}, " API base paths")

	_ = UpdateCmd.MarkFlagRequired("id")
	_ = UpdateCmd.MarkFlagRequired("display-name")
	_ = UpdateCmd.MarkFlagRequired("resource-uri")
	_ = UpdateCmd.MarkFlagRequired("endpoints")
}

// Copyright 2023 Google LLC
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

package specs

import (
	"internal/apiclient"

	"internal/client/registry/apis/versions"

	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// DeleteSpecCmd to get instance
var DeleteSpecCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a spec associated with an API version",
	Long:  "Delete a spec associated with an API version in Apigee Registry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = versions.DeleteSpec(apiName, apiVersion, name)
		return
	},
}

func init() {
	DeleteSpecCmd.Flags().StringVarP(&apiName, "api-name", "",
		"", "API Name")
	DeleteSpecCmd.Flags().StringVarP(&apiVersion, "api-version", "",
		"", "API Version")
	DeleteSpecCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Spec")

	_ = DeleteSpecCmd.MarkFlagRequired("api-name")
	_ = DeleteSpecCmd.MarkFlagRequired("api-version")
	_ = DeleteSpecCmd.MarkFlagRequired("name")
}

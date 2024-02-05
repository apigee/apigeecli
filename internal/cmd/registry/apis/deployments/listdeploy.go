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

package deployments

import (
	"internal/apiclient"

	"internal/client/registry/apis"

	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// ListDeployCmd to get instance
var ListDeployCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments for an API in Apigee Registry",
	Long:  "List deployments for an API in Apigee Registry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if name != "" {
			_, err = apis.ListDeploymentRevisions(apiName, name, pageSize, pageToken, filter, orderBy)
			return err
		}
		_, err = apis.ListDeployments(apiName, pageSize, pageToken, filter, orderBy)
		return
	},
}

var (
	apiName, pageToken, filter, orderBy string
	pageSize                            int
)

func init() {
	ListDeployCmd.Flags().StringVarP(&apiName, "api-name", "",
		"", "API Name")
	ListDeployCmd.Flags().StringVarP(&name, "name", "n",
		"", "Deployment Name")
	ListDeployCmd.Flags().StringVarP(&pageToken, "page-token", "",
		"", "A page token, received from a previous list call")
	ListDeployCmd.Flags().StringVarP(&filter, "filter", "",
		"", "An expression that can be used to filter the list")
	ListDeployCmd.Flags().StringVarP(&orderBy, "order-by", "",
		"", "A comma-separated list of fields to be sorted; ex: foo desc")

	_ = ListDeployCmd.MarkFlagRequired("api-name")
}

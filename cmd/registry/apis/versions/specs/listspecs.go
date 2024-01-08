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

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// ListSpecCmd to get instance
var ListSpecCmd = &cobra.Command{
	Use:   "list",
	Short: "List all specs for an API version",
	Long:  "List all artifacts for an API version in Apigee Registry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = versions.ListSpecs(apiName, apiVersion, pageSize, pageToken, filter, orderBy)
		return
	},
}

var (
	apiName, apiVersion, pageToken, filter, orderBy string
	pageSize                                        int
)

func init() {
	ListSpecCmd.Flags().StringVarP(&apiName, "api-name", "",
		"", "API Name")
	ListSpecCmd.Flags().StringVarP(&apiVersion, "api-version", "",
		"", "API Version")
	ListSpecCmd.Flags().StringVarP(&pageToken, "page-token", "",
		"", "A page token, received from a previous list call")
	ListSpecCmd.Flags().StringVarP(&filter, "filter", "",
		"", "An expression that can be used to filter the list")
	ListSpecCmd.Flags().StringVarP(&orderBy, "order-by", "",
		"", "A comma-separated list of fields to be sorted; ex: foo desc")

	_ = ListSpecCmd.MarkFlagRequired("api-name")
	_ = ListSpecCmd.MarkFlagRequired("api-version")
}

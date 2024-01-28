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

package versions

import (
	"internal/apiclient"

	"internal/client/registry/apis"

	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// ListVersionCmd to get instance
var ListVersionCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API Versions for an API",
	Long:  "List all API Versions for an API in Apigee Registry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.ListVersions(name, pageSize, pageToken, filter, orderBy)
		return
	},
}

var (
	name, pageToken, filter, orderBy string
	pageSize                         int
)

func init() {
	ListVersionCmd.Flags().StringVarP(&name, "name", "n",
		"", "Apigee Registry API Name")
	ListVersionCmd.Flags().StringVarP(&pageToken, "page-token", "",
		"", "A page token, received from a previous list call")
	ListVersionCmd.Flags().StringVarP(&filter, "filter", "",
		"", "An expression that can be used to filter the list")
	ListVersionCmd.Flags().StringVarP(&orderBy, "order-by", "",
		"", "A comma-separated list of fields to be sorted; ex: foo desc")

	_ = ListVersionCmd.MarkFlagRequired("name")
}

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

package apps

import (
	"internal/apiclient"

	"internal/client/apps"

	"github.com/spf13/cobra"
)

// ListCmd to list apps
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of Developer Applications",
	Long:  "Returns a list of app IDs within an organization based on app status",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apps.List(includeCred, expand, count, status,
			startKey, ids, keyStatus, apiProduct, pageSize, pageToken, filter)
		return
	},
}

var (
	status, startKey, ids, keyStatus, apiProduct, pageToken, filter string
	expand                                                          = false
	includeCred                                                     = false
	count, pageSize                                                 int
)

func init() {
	ListCmd.Flags().IntVarP(&count, "count", "c",
		-1, "Number of apps; limit is 1000")
	ListCmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "Expand Details")
	ListCmd.Flags().BoolVarP(&includeCred, "incl-cred", "i",
		false, "Include Credentials")
	ListCmd.Flags().StringVarP(&status, "status", "s",
		"", "Filter by the status of the app. Valid values are approved or revoked")
	ListCmd.Flags().StringVarP(&ids, "ids", "",
		"", "Comma-separated list of app IDs")
	ListCmd.Flags().StringVarP(&keyStatus, "key-status", "k",
		"", "Key status of the app. Valid values include approved or revoked")
	ListCmd.Flags().StringVarP(&apiProduct, "api-product", "p",
		"", "Name of the API Product to filter by")
	ListCmd.Flags().IntVarP(&pageSize, "page-size", "",
		-1, "Count of apps a single page can have in the response")
	ListCmd.Flags().StringVarP(&pageToken, "page-token", "",
		"", "The starting index record for listing the apps")
	ListCmd.Flags().StringVarP(&filter, "filter", "f",
		"", "The filter expression to be used to get the list of apps")
}

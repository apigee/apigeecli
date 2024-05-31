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

package products

import (
	"fmt"

	"internal/apiclient"

	"internal/client/products"

	"github.com/spf13/cobra"
)

// ListCmd to list products
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of API products",
	Long:  "Returns a list of API products with a filter by attribute names and values if provided",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if len(filter) > 1 {
			return fmt.Errorf("invalid filter options. Only proxies are supported")
		}
		if len(filter) == 1 && (filter["proxy"] == "") {
			return fmt.Errorf("invalid filter options. Filter option must be proxies, " +
				"expand must be set to true and count cannot be set")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if len(filter) > 0 {
			_, err = products.ListFilter(filter)
		} else {
			_, err = products.List(count, startKey, expand)
		}
		return err
	},
}

var (
	expand   = false
	count    int
	filter   map[string]string
	startKey string
)

func init() {
	ListCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	ListCmd.Flags().IntVarP(&count, "count", "c",
		-1, "Number of products; limit is 1000")

	ListCmd.Flags().StringVarP(&startKey, "start", "s",
		"", "Set an API Product name as the start key")

	ListCmd.Flags().StringToStringVarP(&filter, "filter", "",
		nil, "Add a filter condition. Ex: proxies=foo will return products that contain the proxy foo")

	ListCmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "Expand Details")
}

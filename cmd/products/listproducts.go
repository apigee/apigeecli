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
	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/products"
	"github.com/spf13/cobra"
)

//Cmd to list products
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of API products",
	Long:  "Returns a list of API products with a filter by attribute names and values if provided",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = products.List(count, expand)
		return

	},
}

var expand = false
var count int

func init() {

	ListCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	ListCmd.Flags().IntVarP(&count, "count", "c",
		-1, "Number of products; limit is 1000")

	ListCmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "Expand Details")

}

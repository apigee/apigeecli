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
	"internal/apiclient"
	"internal/client/products"

	"github.com/spf13/cobra"
)

// DelRatePlanCmd to list envs
var DelRatePlanCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete rate plan associated with the API Product",
	Long:  "Delete rate plan associated with the API Product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = products.DeleteRatePlan(apiproduct, rateplan)
		return
	},
}

func init() {
	DelRatePlanCmd.Flags().StringVarP(&apiproduct, "product", "p",
		"", "name of the API Product")
	DelRatePlanCmd.Flags().StringVarP(&rateplan, "rateplan", "",
		"", "Rate Plan Id")

	_ = DelRatePlanCmd.MarkFlagRequired("product")
	_ = DelRatePlanCmd.MarkFlagRequired("rateplan")
}

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

// GetRatePlanCmd to list envs
var GetRatePlanCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a rate plan associated with an API Product",
	Long:  "Get a rate plan associated with an API Product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = products.GetRatePlan(name, rateplan)
		return
	},
}

func init() {
	GetRatePlanCmd.Flags().StringVarP(&name, "name", "n",
		"", "name of the API Product")
	GetRatePlanCmd.Flags().StringVarP(&rateplan, "rateplan", "p",
		"", "name of the API Product")

	_ = GetRatePlanCmd.MarkFlagRequired("name")
	_ = GetRatePlanCmd.MarkFlagRequired("rateplan")
}

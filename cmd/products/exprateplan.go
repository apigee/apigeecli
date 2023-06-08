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

// ExpRateplanCmd to export products
var ExpRateplanCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the rate plans of an API product to a file",
	Long:  "Export the rate plans of an API product to a file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		exportFileName := "rateplan_" + name + ".json"
		respBody, err := products.ExportRateplan(name)
		if err != nil {
			return err
		}
		return apiclient.WriteByteArrayToFile(exportFileName, false, respBody)
	},
}

func init() {
	ExpRateplanCmd.Flags().StringVarP(&name, "name", "n",
		"", "name of the API Product")

	_ = ExpRateplanCmd.MarkFlagRequired("name")
}

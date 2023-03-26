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
	"os"

	"internal/apiclient"

	"internal/clilog"

	"github.com/apigee/apigeecli/client/products"
	"github.com/spf13/cobra"
)

// Cmd to create a rate plane for an api product
var CreateRateplanCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a rate plan for an API product",
	Long:  "Create a rate plan for an API product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if rateplanData, err = os.ReadFile(rateplanFile); err != nil {
			clilog.Info.Println(err)
			return
		}
		_, err = products.CreateRatePlan(apiproduct, rateplanData)
		return
	},
}

var apiproduct, rateplanFile string
var rateplanData []byte

func init() {

	CreateRateplanCmd.Flags().StringVarP(&apiproduct, "product", "p",
		"", "Name of the API Product")
	CreateRateplanCmd.Flags().StringVarP(&rateplanFile, "rateplan", "",
		"", "File containing Rate plane JSON. See samples for how to create the file")

	_ = CreateRateplanCmd.MarkFlagRequired("apiproduct")
	_ = CreateRateplanCmd.MarkFlagRequired("rateplan")
}

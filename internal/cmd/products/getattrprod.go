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

// GetAttrCmd to delete product attribute
var GetAttrCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an attribute of an API product",
	Long:  "Get an attribute of an API product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = products.GetAttribute(name, attrName)
		return
	},
}

func init() {
	GetAttrCmd.Flags().StringVarP(&attrName, "attr", "k",
		"", "API Product attribute name")

	_ = GetAttrCmd.MarkFlagRequired("attr")
}

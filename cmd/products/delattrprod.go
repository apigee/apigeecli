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

// DelAttrCmd to delete product attribute
var DelAttrCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an attribute from an API product",
	Long:  "Deletes an attribute from an API product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = products.DeleteAttribute(name, attrName)
		return
	},
}

func init() {

	DelAttrCmd.Flags().StringVarP(&attrName, "attr", "k",
		"", "API Product attribute name")

	_ = DelAttrCmd.MarkFlagRequired("attr")
}

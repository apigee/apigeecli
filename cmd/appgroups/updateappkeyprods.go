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

package appgroups

import (
	"internal/apiclient"
	"internal/client/appgroups"

	"github.com/spf13/cobra"
)

// UpdateKeyProdCmd to create developer keys
var UpdateKeyProdCmd = &cobra.Command{
	Use:   "update-prod",
	Short: "Update products in an app key contained in an AppGroup",
	Long:  "Update products in an app key contained in an AppGroup",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = appgroups.UpdateKeyProducts(name, appName, key, apiProducts)
		return
	},
}

func init() {
	UpdateKeyProdCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the AppGroup")
	UpdateKeyProdCmd.Flags().StringVarP(&appName, "app-name", "",
		"", "Name of the app")
	UpdateKeyProdCmd.Flags().StringVarP(&key, "key", "k",
		"", "AppGroup app consumer key")
	UpdateKeyProdCmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")

	_ = UpdateKeyProdCmd.MarkFlagRequired("name")
	_ = UpdateKeyProdCmd.MarkFlagRequired("app-name")
	_ = UpdateKeyProdCmd.MarkFlagRequired("key")
	_ = UpdateKeyProdCmd.MarkFlagRequired("prods")
}

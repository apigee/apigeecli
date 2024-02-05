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

// ImpCmd to import products
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing API products",
	Long:  "Import a file containing API products",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return products.Import(conn, filePath, upsert)
	},
}

var (
	upsert   bool
	filePath string
)

func init() {
	ImpCmd.Flags().StringVarP(&filePath, "file", "f",
		"", "File containing API Products")
	ImpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
	ImpCmd.Flags().BoolVarP(&upsert, "upsert", "",
		false, "Insert or update products")

	_ = ImpCmd.MarkFlagRequired("file")
}

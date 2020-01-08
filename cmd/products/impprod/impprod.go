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

package impprod

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/client/products"
)

//Cmd to import products
var Cmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing API products",
	Long:  "Import a file containing API products",
	RunE: func(cmd *cobra.Command, args []string) error {
		return products.Import(conn, filePath)
	},
}

var conn int
var filePath string

func init() {

	Cmd.Flags().StringVarP(&filePath, "file", "f",
		"", "File containing API Products")
	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = Cmd.MarkFlagRequired("file")
}

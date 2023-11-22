// Copyright 2023 Google LLC
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

package apicategories

import (
	"internal/apiclient"

	"internal/client/apicategories"

	"github.com/spf13/cobra"
)

// CreateCmd to get a catalog items
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new API category",
	Long:  "Creates a new API category",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apicategories.Create(siteid, name)
		return
	},
}

var name string

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Category name")
	_ = CreateCmd.MarkFlagRequired("name")
}

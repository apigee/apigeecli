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

// MoveCmd to move products between spaces
var MoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move an API product between Spaces",
	Long:  "Move an API product between Spaces",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = products.Move(name, space)
		return err
	},
}

func init() {
	MoveCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the API Product")
	MoveCmd.Flags().StringVarP(&space, "space", "",
		"", "Name of the Apigee Space moving to")

	// TODO: apiresource -r later

	_ = UpdateCmd.MarkFlagRequired("name")
	_ = UpdateCmd.MarkFlagRequired("space")
}

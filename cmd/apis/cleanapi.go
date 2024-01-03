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

package apis

import (
	"internal/apiclient"

	"internal/client/apis"

	"github.com/spf13/cobra"
)

// CleanCmd to delete api
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Deletes undeployed/unused revisions of an API proxy",
	Long:  "Deletes undeployed/unused revisions of an API proxy",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return apis.CleanProxy(name, reportOnly, keepList)
	},
}

var (
	reportOnly bool
	keepList   []string
)

func init() {
	CleanCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	CleanCmd.Flags().BoolVarP(&reportOnly, "report", "",
		true, "Report which API proxy revisions will be deleted")
	CleanCmd.Flags().StringArrayVarP(&keepList, "keepList", "k",
		[]string{}, "List of revisions to keep, -k 1 -k 2")

	_ = CleanCmd.MarkFlagRequired("name")
}

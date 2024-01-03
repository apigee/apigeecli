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

package appgroups

import (
	"internal/apiclient"

	"internal/client/appgroups"

	"github.com/spf13/cobra"
)

// ListCmd to list apps
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of AppGroups",
	Long:  "Returns a list of AppGroups",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = appgroups.List(pageSize, pageToken, filter)
		return
	},
}

var (
	filter    string
	pageToken string
	pageSize  int
)

func init() {
	ListCmd.Flags().IntVarP(&pageSize, "page-size", "s",
		-1, "Number of appgroups; limit is 1000")
	ListCmd.Flags().StringVarP(&pageToken, "page-token", "p",
		"", "Page token")
	ListCmd.Flags().StringVarP(&filter, "filter", "f",
		"", "List filter")
}

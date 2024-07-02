// Copyright 2024 Google LLC
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

package versions

import (
	"internal/apiclient"
	"internal/client/hub"

	"github.com/spf13/cobra"
)

// ListCmd to get a catalog items
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API Hub API Versions",
	Long:  "List API Hub API Versions",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = hub.ListApiVersions(id, filter, pageSize, pageToken)
		return
	},
}

var (
	pageSize              int
	pageToken, filter, id string
)

func init() {
	ListCmd.Flags().StringVarP(&id, "api-id", "i",
		"", "API ID")
	ListCmd.Flags().StringVarP(&filter, "filter", "f",
		"", "filter expression")
	ListCmd.Flags().IntVarP(&pageSize, "page-size", "",
		-1, "The maximum number of versions to return")
	ListCmd.Flags().StringVarP(&pageToken, "page-token", "",
		"", "A page token, received from a previous call")

	_ = ListCmd.MarkFlagRequired("api-id")
}

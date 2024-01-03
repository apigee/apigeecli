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

package developers

import (
	"internal/apiclient"

	"internal/client/developers"

	"github.com/spf13/cobra"
)

// ListCmd to list developer
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of App Developers",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	Long: "Lists all developers in an organization by email address",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = developers.List(count, expand, ids)
		return
	},
}

var (
	count int
	ids   string
)

func init() {
	ListCmd.Flags().IntVarP(&count, "count", "c",
		-1, "Number of developers; limit is 1000")

	ListCmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "Expand Details")

	ListCmd.Flags().StringVarP(&ids, "ids", "i",
		"", "List of IDs to include, separated by commas")
}

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
		cmd.SilenceUsage = true

		_, err = developers.List(count, expand, ids, startKey, app, includeCompany)
		return
	},
}

var (
	count              int
	ids, startKey, app string
	includeCompany     bool
)

func init() {
	ListCmd.Flags().IntVarP(&count, "count", "c",
		-1, "Number of developers; limit is 1000")

	ListCmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "Expand Details")

	ListCmd.Flags().StringVarP(&ids, "ids", "i",
		"", "List of IDs to include, separated by commas")

	ListCmd.Flags().StringVarP(&startKey, "start-key", "s",
		"", "Email address of the developer from which to start displaying the list of developers")

	ListCmd.Flags().BoolVarP(&includeCompany, "include-company", "",
		false, "Flag that specifies whether to include company details in the response")

	ListCmd.Flags().StringVarP(&app, "app", "",
		"", "List only Developers that are associated with the app")
}

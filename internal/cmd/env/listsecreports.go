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

package env

import (
	"internal/apiclient"
	"internal/client/env"

	"github.com/spf13/cobra"
)

// ListSecReportCmd returns security incidents
var ListSecReportCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns a security reports by name",
	Long:  "Returns a security reports by name",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = env.ListSecurityReports(pageSize, pageToken, dataset, to,
			from, status, submittedBy)
		return
	},
}

var from, to, dataset, status, submittedBy string

func init() {
	ListSecReportCmd.Flags().IntVarP(&pageSize, "page-size", "",
		-1, "The maximum number of versions to return")
	ListSecReportCmd.Flags().StringVarP(&pageToken, "page-token", "",
		"", "A page token, received from a previous call")
	ListSecReportCmd.Flags().StringVarP(&from, "from", "f",
		"", "Filter response list by returning security reports that created after this date time")
	ListSecReportCmd.Flags().StringVarP(&to, "to", "",
		"", "Filter response list by returning security reports that created before this date time")
	ListSecReportCmd.Flags().StringVarP(&dataset, "dataset", "d",
		"", "Filter response list by dataset; example api, mint")
	ListSecReportCmd.Flags().StringVarP(&status, "status", "s",
		"", "Filter response list by security reports status")
	ListSecReportCmd.Flags().StringVarP(&submittedBy, "submitted-by", "u",
		"", "Filter response list by user who submitted queries")
}

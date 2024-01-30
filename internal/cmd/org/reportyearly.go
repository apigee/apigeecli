// Copyright 2022 Google LLC
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

package org

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"internal/apiclient"
	"internal/clilog"

	"internal/client/orgs"

	"github.com/spf13/cobra"
)

// YearlyCmd to get monthly usage
var YearlyCmd = &cobra.Command{
	Use:   "yearly",
	Short: "Report yearly usage for an Apigee Org",
	Long:  "Report yearly usage for an Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var apiCalls int
		w := tabwriter.NewWriter(os.Stdout, 32, 4, 0, ' ', 0)

		clilog.Warning.Println("This API is rate limited to 1 API Call per second")

		if _, err = time.Parse("2006", fmt.Sprintf("%d", year)); err != nil {
			return err
		}

		if envDetails {
			fmt.Fprintln(w, "ENVIRONMENT\tMONTH\tAPI CALLS")
			w.Flush()
		}

		if apiCalls, err = orgs.TotalAPICallsInYear(year, envDetails, conn); err != nil {
			return err
		}

		if envDetails {
			fmt.Printf("\nSummary\n\n")
		}

		fmt.Fprintln(w, "ORGANIZATION\tYEAR\tAPI CALLS")
		fmt.Fprintf(w, "%s\t%d\t%d\n", apiclient.GetApigeeOrg(), year, apiCalls)
		fmt.Fprintln(w)
		w.Flush()

		return err
	},
}

func init() {
	YearlyCmd.Flags().IntVarP(&year, "year", "y",
		-1, "Year")
	YearlyCmd.Flags().BoolVarP(&envDetails, "env-details", "",
		false, "Print details of each environment")
	YearlyCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
	YearlyCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	_ = YearlyCmd.MarkFlagRequired("year")
}

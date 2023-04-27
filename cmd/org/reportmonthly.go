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

// MonthlyCmd to get monthly usage
var MonthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Report monthly usage for an Apigee Org",
	Long:  "Report monthly usage for an Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var apiCalls int
		w := tabwriter.NewWriter(os.Stdout, 32, 4, 0, ' ', 0)

		clilog.Warning.Println("This API is rate limited to 1 API Call per second")

		if _, err = time.Parse("1/2006", fmt.Sprintf("%d/%d", month, year)); err != nil {
			return
		}

		if envDetails {
			fmt.Fprintln(w, "ENVIRONMENT\tMONTH\tAPI CALLS")
			w.Flush()
		}

		if apiCalls, err = orgs.TotalAPICallsInMonth(month, year, envDetails, conn); err != nil {
			return
		}

		if envDetails {
			fmt.Printf("\nSummary\n\n")
		}

		fmt.Fprintln(w, "ORGANIATION\tMONTH\tAPI CALLS")
		fmt.Fprintf(w, "%s\t%d/%d\t%d\n", apiclient.GetApigeeOrg(), month, year, apiCalls)
		fmt.Fprintln(w)
		w.Flush()

		return
	},
}

var (
	month, year int
	envDetails  bool
)

func init() {
	MonthlyCmd.Flags().IntVarP(&month, "month", "m",
		-1, "Month")
	MonthlyCmd.Flags().IntVarP(&year, "year", "y",
		-1, "Year")
	MonthlyCmd.Flags().BoolVarP(&envDetails, "env-details", "",
		false, "Print details of each environment")
	MonthlyCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
	MonthlyCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	_ = MonthlyCmd.MarkFlagRequired("month")
	_ = MonthlyCmd.MarkFlagRequired("year")
}

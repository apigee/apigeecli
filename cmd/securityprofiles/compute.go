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

package securityprofiles

import (
	"fmt"
	"time"

	"internal/apiclient"
	"internal/client/securityprofiles"

	"github.com/spf13/cobra"
)

// ComputeCmd to list catalog items
var ComputeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Calculates scores for requested time range",
	Long:  "Calculates scores for requested time range for the specified security profile",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if startTime != "" {
			if _, err = time.Parse(time.RFC3339, startTime); err != nil {
				return fmt.Errorf("invalid format for startTime: %v", err)
			}
		}
		if endTime != "" {
			if _, err = time.Parse(time.RFC3339, endTime); err != nil {
				return fmt.Errorf("invalid format for endTime: %v", err)
			}
		}
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if endTime == "" {
			// if endTime is not set, then set current timestamp
			endTime = time.Now().UTC().Format(time.RFC3339)
		}
		if startTime == "" {
			// if startTime is not set, then set yesterday's timestamp
			startTime = time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339)
		}
		_, err = securityprofiles.Compute(name, startTime, endTime,
			filters, pageSize, pageToken, full)
		return
	},
}

var (
	startTime, endTime string
	filters            []string
	full               bool
)

func init() {
	ComputeCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the security profile")
	ComputeCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")

	ComputeCmd.Flags().StringVarP(&startTime, "start-time", "",
		"", "Inclusive start of the interval; default is 24 hours ago")
	ComputeCmd.Flags().StringVarP(&endTime, "end-time", "",
		"", "Exclusive end of the interval; default current timestamp")

	ComputeCmd.Flags().StringArrayVarP(&filters, "filters", "f",
		nil, "Filters are used to filter scored components")
	ComputeCmd.Flags().IntVarP(&pageSize, "page-size", "",
		-1, "The maximum number of versions to return")
	ComputeCmd.Flags().StringVarP(&pageToken, "page-token", "",
		"", "A page token, received from a previous call")
	ComputeCmd.Flags().BoolVarP(&full, "full", "",
		false, "Generate all scores")

	_ = ComputeCmd.MarkFlagRequired("name")
	_ = ComputeCmd.MarkFlagRequired("env")
}

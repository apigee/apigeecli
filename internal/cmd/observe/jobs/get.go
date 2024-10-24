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

package jobs

import (
	"internal/apiclient"
	"internal/client/observe"

	"github.com/spf13/cobra"
)

// GetCmd to get a catalog items
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details for an Observation Job",
	Long:  "Get details for an Observation Job",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = observe.GetObservationJob(observationJobId)
		return
	},
}

func init() {
	GetCmd.Flags().StringVarP(&observationJobId, "id", "i",
		"", "Observation Job Id")

	_ = GetCmd.MarkFlagRequired("id")
}

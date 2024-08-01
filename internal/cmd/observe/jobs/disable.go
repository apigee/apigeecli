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

// DisableCmd to get a catalog items
var DisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable an Observation Job",
	Long:  "Disable an Observation Job",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = observe.DisableObservationJob(observationJobId)
		return
	},
}

func init() {
	DisableCmd.Flags().StringVarP(&observationJobId, "id", "i",
		"", "Observation Job Id")

	_ = DisableCmd.MarkFlagRequired("id")
}

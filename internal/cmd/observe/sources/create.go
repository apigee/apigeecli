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

package sources

import (
	"internal/apiclient"
	"internal/client/observe"

	"github.com/spf13/cobra"
)

// CrtCmd
var CrtCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Observation Source",
	Long:  "Create a new Observation Source",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = observe.CreateObservationSource(observationSourceId, pscConfig)
		return
	},
}

var (
	observationSourceId string
	pscConfig           map[string]string
)

func init() {
	CrtCmd.Flags().StringVarP(&observationSourceId, "id", "i",
		"", "Observation Job")
	CrtCmd.Flags().StringToStringVarP(&pscConfig, "psc-config", "p",
		nil, "The VPC networks where traffic will be observed")

	_ = CrtCmd.MarkFlagRequired("id")
	_ = CrtCmd.MarkFlagRequired("sources")
}

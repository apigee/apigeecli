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

package datacollectors

import (
	"internal/apiclient"

	"internal/client/datacollectors"

	"github.com/spf13/cobra"
)

// ExpCmd to create a new data collector
var ExpCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Data Collectors",
	Long:  "Export Data Collectors",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		const dataCollFileName = "datacollectors.json"
		respBody, err := datacollectors.List()
		if err != nil {
			return err
		}
		return apiclient.WriteByteArrayToFile(dataCollFileName, false, respBody)
	},
}

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

package sharedflows

import (
	"internal/apiclient"

	"internal/client/sharedflows"

	"github.com/spf13/cobra"
)

// CleanCmd to delete sf
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Deletes undeployed/unused reivisions of a Sharedflow",
	Long:  "Deletes undeployed/unused reivisions of a Sharedflow",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		return sharedflows.Clean(name, reportOnly)
	},
}

var reportOnly bool

func init() {
	CleanCmd.Flags().StringVarP(&name, "name", "n",
		"", "Sharedflow name")
	CleanCmd.Flags().BoolVarP(&reportOnly, "report", "",
		true, "Report which Shareflow revisions will be deleted")

	_ = CleanCmd.MarkFlagRequired("name")
}

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

package reports

import (
	"internal/apiclient"
	"internal/client/reports"

	"github.com/spf13/cobra"
)

// DelCmd to get a resource
var DelCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a custom report",
	Long:  "Delete a custom report",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = reports.Delete(name)
		return err
	},
}

func init() {
	DelCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the custom report")

	_ = DelCmd.MarkFlagRequired("name")
}

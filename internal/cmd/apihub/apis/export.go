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

package apis

import (
	"internal/apiclient"
	"internal/client/hub"
	"os"

	"github.com/spf13/cobra"
)

// ExportCmd to get a catalog items
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export API, versions and specifications",
	Long:  "Export API, versions and specifications",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		if folder == "" {
			folder, _ = os.Getwd()
		}
		err = hub.ExportApi(apiID, folder)
		return
	},
}

var folder string

func init() {
	ExportCmd.Flags().StringVarP(&apiID, "api-id", "",
		"", "API ID")
	ExportCmd.Flags().StringVarP(&folder, "folder", "",
		"", "Folder to export the API details")

	_ = ExportCmd.MarkFlagRequired("api-id")
}

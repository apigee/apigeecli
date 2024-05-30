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

package apicategories

import (
	"os"

	"internal/apiclient"

	"internal/client/apicategories"

	"github.com/spf13/cobra"
)

// ExpCmd to export apidocs
var ExpCmd = &cobra.Command{
	Use:   "export",
	Short: "Export API Categories across all sites",
	Long:  "Export API Categories across all sites",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if folder == "" {
			folder, _ = os.Getwd()
		}
		if err = apiclient.FolderExists(folder); err != nil {
			return err
		}
		apiclient.DisableCmdPrintHttpResponse()
		return apicategories.Export(folder)
	},
}

var folder string

func init() {
	ExpCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "folder to export API Docs")
}

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

package appgroups

import (
	"fmt"

	"internal/apiclient"

	"internal/client/appgroups"

	"github.com/spf13/cobra"
)

// ImpAppCmd to import apps
var ImpAppCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing Apps to an AppGroup",
	Long:  "Import a file containing Apps to an AppGroup",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if name == "" && !all {
			return fmt.Errorf("either all must be set to true or a name must be passed")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		apiclient.DisableCmdPrintHttpResponse()
		if name != "" {
			return appgroups.ImportApps(conn, filePath, name)
		}
		return appgroups.ImportAllApps(conn, filePath)
	},
}

func init() {
	ImpAppCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the AppGroup")
	ImpAppCmd.Flags().StringVarP(&filePath, "file", "f",
		"", "File containing Apps")
	ImpAppCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
	ImpAppCmd.Flags().BoolVarP(&all, "all", "",
		false, "Import Apps for all AppGroups in the org")

	_ = ImpAppCmd.MarkFlagRequired("file")
}

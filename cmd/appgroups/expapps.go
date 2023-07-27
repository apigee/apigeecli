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

package appgroups

import (
	"fmt"

	"internal/apiclient"

	"internal/client/appgroups"

	"github.com/spf13/cobra"
)

// ExpAppCmd to export apps
var ExpAppCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Apps in an AppGroup to a file",
	Long:  "Export Apps in an AppGroup to a file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if name == "" && !all {
			return fmt.Errorf("either all must be set to true or a name must be passed")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		exportFileName := name + "_" + "apps.json"

		apiclient.DisableCmdPrintHttpResponse()

		if name != "" {
			payload, err := appgroups.ExportApps(name)
			if err != nil {
				return err
			}
			prettyPayload, err := apiclient.PrettifyJSON(payload)
			if err != nil {
				return err
			}
			return apiclient.WriteByteArrayToFile(exportFileName, false, prettyPayload)
		}

		appGroupsList, err := appgroups.Export()
		if err != nil {
			return err
		}

		payload, err := appgroups.ExportAllApps(appGroupsList, 4)
		if err != nil {
			return err
		}
		return apiclient.WriteArrayByteArrayToFile(exportFileName, false, payload)
	},
}

var all bool

func init() {
	ExpAppCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the AppGroup")
	ExpAppCmd.Flags().BoolVarP(&all, "all", "",
		false, "Export apps for all appgroups in the org")
}

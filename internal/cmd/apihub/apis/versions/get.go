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

package versions

import (
	"internal/apiclient"
	"internal/client/hub"

	"github.com/spf13/cobra"
)

// GetCmd to get a catalog items
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details for an API Version",
	Long:  "Get details for an API Version",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		if definition != "" {
			_, err = hub.GetApiVersionsDefinitions(versionID, apiID, definition)
			return
		}
		_, err = hub.GetApiVersion(versionID, apiID)
		return
	},
}

var versionID, definition string

func init() {
	GetCmd.Flags().StringVarP(&versionID, "id", "i",
		"", "API Version ID")
	GetCmd.Flags().StringVarP(&apiID, "api-id", "",
		"", "API ID")
	GetCmd.Flags().StringVarP(&definition, "definition", "d",
		"", "Get API Version definition")

	_ = GetCmd.MarkFlagRequired("api-id")
	_ = GetCmd.MarkFlagRequired("version")
}

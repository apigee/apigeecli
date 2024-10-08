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
	"os"

	"github.com/spf13/cobra"
)

// UpdateCmd
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an API Hub API Version",
	Long:  "Update an API Hub API Version",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		var apiFileBytes []byte

		if apiFileBytes, err = os.ReadFile(apiFilePath); err != nil {
			return err
		}
		_, err = hub.UpdateApiVersion(id, apiID, apiFileBytes)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&id, "id", "i",
		"", "API Version ID")
	UpdateCmd.Flags().StringVarP(&apiID, "api-id", "",
		"", "API ID")
	UpdateCmd.Flags().StringVarP(&apiFilePath, "file", "f",
		"", "Path to a file containing the API version definition")

	_ = UpdateCmd.MarkFlagRequired("api-id")
	_ = UpdateCmd.MarkFlagRequired("file")
}

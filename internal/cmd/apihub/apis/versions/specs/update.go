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

package specs

import (
	"path/filepath"

	"internal/apiclient"
	"internal/client/hub"
	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// UpdateCmd
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates a Spec for an API Hub API Version",
	Long:  "Updates a Spec for an API Hub API Version",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		var contents []byte

		if contents, err = utils.ReadFile(apiFilePath); err != nil {
			return err
		}
		_, err = hub.UpdateApiVersionSpec(apiID, versionID, specID, displayName,
			contents, filepath.Ext(apiFilePath), sourceURI)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&specID, "id", "i",
		"", "Spec ID")
	UpdateCmd.Flags().StringVarP(&versionID, "version", "v",
		"", "API Version ID")
	UpdateCmd.Flags().StringVarP(&apiID, "api-id", "",
		"", "API ID")
	UpdateCmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "Spec Display Name")
	UpdateCmd.Flags().StringToStringVar(&attributes, "attrs",
		nil, "API Spec attributes")
	UpdateCmd.Flags().StringVarP(&sourceURI, "source-uri", "s",
		"", "API Spec attributes")
	UpdateCmd.Flags().StringVarP(&apiFilePath, "file", "f",
		"", "Path to a file containing the API spec")

	_ = UpdateCmd.MarkFlagRequired("id")
	_ = UpdateCmd.MarkFlagRequired("api-id")
	_ = UpdateCmd.MarkFlagRequired("version")
	_ = UpdateCmd.MarkFlagRequired("display-name")
	_ = UpdateCmd.MarkFlagRequired("file")
}

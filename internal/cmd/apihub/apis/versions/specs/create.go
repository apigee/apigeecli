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
	"internal/apiclient"
	"internal/client/hub"
	"internal/cmd/utils"
	"path/filepath"

	"github.com/spf13/cobra"
)

// CrtCmd
var CrtCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Spec for an API Hub API Version",
	Long:  "Create a new Spec for an API Hub API Version",
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
		_, err = hub.CreateApiVersionsSpec(apiID, versionID, specID, displayName,
			contents, filepath.Ext(apiFilePath), sourceURI, documentation)
		return
	},
	Example: `Create a new API Version Specification: ` + GetExample(0),
}

var (
	specID, displayName, mimeType, sourceURI, documentation, apiFilePath string
	attributes                                                           map[string]string
)

func init() {
	CrtCmd.Flags().StringVarP(&specID, "id", "i",
		"", "Spec ID")
	CrtCmd.Flags().StringVarP(&versionID, "version", "v",
		"", "API Version ID")
	CrtCmd.Flags().StringVarP(&apiID, "api-id", "",
		"", "API ID")
	CrtCmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "Spec Display Name")
	CrtCmd.Flags().StringToStringVar(&attributes, "attrs",
		nil, "API Spec attributes")
	CrtCmd.Flags().StringVarP(&sourceURI, "source-uri", "s",
		"", "API Spec attributes")
	CrtCmd.Flags().StringVarP(&documentation, "documentation", "",
		"", "API Spec external documentation")
	CrtCmd.Flags().StringVarP(&apiFilePath, "file", "f",
		"", "Path to a file containing the API spec")

	_ = CrtCmd.MarkFlagRequired("api-id")
	_ = CrtCmd.MarkFlagRequired("version")
	_ = CrtCmd.MarkFlagRequired("display-name")
	_ = CrtCmd.MarkFlagRequired("file")
}

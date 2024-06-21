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

	"github.com/spf13/cobra"
)

// GetCmd to get a catalog items
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a spec for an API Version",
	Long:  "Get a spec for an API Version",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		if contents {
			_, err = hub.GetApiVersionsSpecContents(apiID, versionID, specID)
			return
		}
		_, err = hub.GetApiVersionSpec(apiID, versionID, specID)
		return
	},
}

var contents bool

func init() {
	GetCmd.Flags().StringVarP(&versionID, "version", "v",
		"", "API Version ID")
	GetCmd.Flags().StringVarP(&apiID, "api-id", "",
		"", "API ID")
	GetCmd.Flags().StringVarP(&specID, "spec-id", "s",
		"", "Spec ID")
	GetCmd.Flags().BoolVarP(&contents, "contents", "c",
		false, "Get contents")

	_ = GetCmd.MarkFlagRequired("api-id")
	_ = GetCmd.MarkFlagRequired("version")
	_ = GetCmd.MarkFlagRequired("spec-id")
}

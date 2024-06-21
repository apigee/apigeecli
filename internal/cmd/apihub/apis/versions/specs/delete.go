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

// DelCmd to get a catalog items
var DelCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a spec from an API Version",
	Long:  "Delete a spec from an API Version",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = hub.DeleteApiVersionSpec(apiID, versionID, specID)
		return
	},
}

func init() {
	DelCmd.Flags().StringVarP(&versionID, "version", "v",
		"", "API Version ID")
	DelCmd.Flags().StringVarP(&apiID, "api-id", "",
		"", "API ID")
	DelCmd.Flags().StringVarP(&specID, "spec-id", "s",
		"", "Spec ID")

	_ = DelCmd.MarkFlagRequired("api-id")
	_ = DelCmd.MarkFlagRequired("version")
	_ = DelCmd.MarkFlagRequired("spec-id")
}

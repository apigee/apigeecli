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

package artifacts

import (
	"internal/apiclient"

	"internal/client/registry/apis/versions"

	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// GetArtifactCmd to get instance
var GetArtifactCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Artifact details for an API version",
	Long:  "Get Artifact details for an API version",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if content {
			return versions.GetArtifactContents(apiName, apiVersion, name)
		}
		_, err = versions.GetArtifact(apiName, apiVersion, name)
		return
	},
}

var content bool

func init() {
	GetArtifactCmd.Flags().StringVarP(&apiName, "api-name", "",
		"", "API Name")
	GetArtifactCmd.Flags().StringVarP(&apiVersion, "api-version", "",
		"", "API Version")
	GetArtifactCmd.Flags().StringVarP(&name, "name", "n",
		"", "Artifact Name")
	GetArtifactCmd.Flags().BoolVarP(&content, "content", "c",
		false, "If set to true, returns artifact contents")

	_ = GetArtifactCmd.MarkFlagRequired("api-name")
	_ = GetArtifactCmd.MarkFlagRequired("api-version")
	_ = GetArtifactCmd.MarkFlagRequired("name")
}

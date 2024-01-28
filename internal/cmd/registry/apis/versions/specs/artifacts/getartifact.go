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
	Short: "Get artifact details for a spec in an API version",
	Long:  "Get artifact details for a spec in an API version in Apigee Registry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = versions.GetSpecArtifact(apiName, apiVersion, specName, name)
		return
	},
}

func init() {
	GetArtifactCmd.Flags().StringVarP(&apiName, "api-name", "",
		"", "API Name")
	GetArtifactCmd.Flags().StringVarP(&apiVersion, "api-version", "",
		"", "API Version")
	GetArtifactCmd.Flags().StringVarP(&specName, "spec-name", "",
		"", "API Version Spec name")
	GetArtifactCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Version Spec Artifact name")

	_ = GetArtifactCmd.MarkFlagRequired("api-name")
	_ = GetArtifactCmd.MarkFlagRequired("api-version")
	_ = GetArtifactCmd.MarkFlagRequired("spec-name")
	_ = GetArtifactCmd.MarkFlagRequired("name")
}

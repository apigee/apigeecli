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

	"internal/client/registry/artifacts"

	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// DelArtifactCmd to get instance
var DelArtifactCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete artifact from Apigee Registry",
	Long:  "Delete artifact from Apigee Registry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = artifacts.Delete(name)
		return
	},
}

func init() {
	DelArtifactCmd.Flags().StringVarP(&name, "name", "n",
		"", "Artifact Name")
}

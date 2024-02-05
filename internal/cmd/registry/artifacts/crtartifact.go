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
	"encoding/base64"

	"internal/apiclient"
	"internal/client/registry/artifacts"

	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// CreateArtifactCmd to create a new instance
var CreateArtifactCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Apigee Registry Artifact",
	Long:  "Create an Apigee Registry Artifact",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var payload []byte
		if payload, err = utils.ReadFile(contentsFile); err != nil {
			return err
		}
		contents := base64.StdEncoding.EncodeToString(payload)
		_, err = artifacts.Create(artifactID, name, contents, labels, annotations)
		return err
	},
}

var name, artifactID, contentsFile string

func init() {
	CreateArtifactCmd.Flags().StringVarP(&artifactID, "id", "i",
		"", "Apigee Registry Artifact ID")
	CreateArtifactCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Artifact")
	CreateArtifactCmd.Flags().StringVarP(&contentsFile, "file", "f",
		"", "Path to a file containing Artifact Contents")
	CreateArtifactCmd.Flags().StringToStringVar(&labels, "labels",
		nil, "Labels attach identifying metadata to resources")
	CreateArtifactCmd.Flags().StringToStringVar(&annotations, "annotations",
		nil, "Annotations attach non-identifying metadata to resources")

	_ = CreateArtifactCmd.MarkFlagRequired("id")
	_ = CreateArtifactCmd.MarkFlagRequired("name")
	_ = CreateArtifactCmd.MarkFlagRequired("file")
}

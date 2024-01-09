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

package specs

import (
	"internal/apiclient"

	"internal/client/registry/apis/versions"

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// CreateSpecCmd to get instance
var CreateSpecCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a spec associated with an API version",
	Long:  "Create a spec associated with an API version in Apigee Registry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fileContents, err := utils.ReadFile(filePath)
		if err != nil {
			return err
		}
		_, err = versions.CreateSpec(apiName, apiVersion, apiSpecId,
			name, fileName, description, sourceURI, string(fileContents), labels, annotations)
		return
	},
}

var labels, annotations map[string]string
var apiSpecId, fileName, description, sourceURI, filePath string

func init() {
	CreateSpecCmd.Flags().StringVarP(&apiName, "api-name", "",
		"", "API Name")
	CreateSpecCmd.Flags().StringVarP(&apiVersion, "api-version", "",
		"", "API Version")
	CreateSpecCmd.Flags().StringVarP(&apiVersion, "api-specid", "",
		"", "API Spec Id")
	CreateSpecCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Version Spec name")
	CreateSpecCmd.Flags().StringVarP(&fileName, "file-name", "",
		"", "A possibly-hierarchical name used to refer to the spec from other specs")
	CreateSpecCmd.Flags().StringVarP(&description, "desc", "d",
		"", "A detailed description for the spec")
	CreateSpecCmd.Flags().StringVarP(&filePath, "file-path", "f",
		"", "A path to file containing the spec")
	CreateSpecCmd.Flags().StringVarP(&sourceURI, "source-uri", "s",
		"", "Source URI")
	CreateSpecCmd.Flags().StringToStringVar(&labels, "labels",
		nil, "Labels attach identifying metadata to resources")
	CreateSpecCmd.Flags().StringToStringVar(&annotations, "annotations",
		nil, "Annotations attach non-identifying metadata to resources")

	_ = CreateSpecCmd.MarkFlagRequired("api-name")
	_ = CreateSpecCmd.MarkFlagRequired("api-version")
	_ = CreateSpecCmd.MarkFlagRequired("name")
	_ = CreateSpecCmd.MarkFlagRequired("file-path")
}

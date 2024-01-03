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

package apis

import (
	"internal/apiclient"
	"internal/client/registry/apis"

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// CreateAPICmd to create a new instance
var CreateAPICmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Apigee Registry API entity",
	Long:  "Create an Apigee Registry API entity",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.Create(apiID, name, displayName, description,
			availability, recommendedVersion, recommendedDeployment,
			labels, annotations)
		return err
	},
}

var apiID, availability, recommendedDeployment, recommendedVersion, description, displayName string

func init() {
	CreateAPICmd.Flags().StringVarP(&apiID, "id", "i",
		"", "Apigee Registry API ID")
	CreateAPICmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the API")
	CreateAPICmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "Human friendly display name of the API")
	CreateAPICmd.Flags().StringVarP(&description, "desc", "",
		"", "Description of the API")
	CreateAPICmd.Flags().StringVarP(&recommendedVersion, "version", "v",
		"", "The recommended version of the API; Format: projects/{project}/"+
			"locations/{location}/apis/{api}/versions/{version}")
	CreateAPICmd.Flags().StringVarP(&recommendedDeployment, "deployment", "",
		"", "The recommended deployment of the API; Format: projects/{project}/"+
			"locations/{location}/apis/{api}/deployments/{deployment}")
	CreateAPICmd.Flags().StringToStringVar(&labels, "labels",
		nil, "Labels attach identifying metadata to resources")
	CreateAPICmd.Flags().StringToStringVar(&annotations, "annotations",
		nil, "Annotations attach non-identifying metadata to resources")

	_ = CreateAPICmd.MarkFlagRequired("id")
	_ = CreateAPICmd.MarkFlagRequired("name")
}

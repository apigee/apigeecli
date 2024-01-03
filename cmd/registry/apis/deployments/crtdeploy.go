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

package deployments

import (
	"internal/apiclient"
	"internal/client/registry/apis"

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// CreateDeployCmd to create a new instance
var CreateDeployCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Apigee Registry API entity",
	Long:  "Create an Apigee Registry API entity",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.CreateDeployment(apiDeploymentID, name, displayName, description,
			apiSpecRevision, endpointURI, externalChannelURI, intendedAudience,
			accessGuidance, labels, annotations)
		return err
	},
}

var (
	apiDeploymentID, apiSpecRevision, endpointURI, externalChannelURI string
	name, description, displayName, intendedAudience, accessGuidance  string
	labels, annotations                                               map[string]string
)

func init() {
	CreateDeployCmd.Flags().StringVarP(&apiDeploymentID, "id", "i",
		"", "Apigee Registry Deployment ID")
	CreateDeployCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Deployment")
	CreateDeployCmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "Human friendly display name of the Deployment")
	CreateDeployCmd.Flags().StringVarP(&description, "desc", "",
		"", "Description of the Deployment")
	CreateDeployCmd.Flags().StringVarP(&apiSpecRevision, "spec-revision", "",
		"", "API Spec Revision")
	CreateDeployCmd.Flags().StringVarP(&endpointURI, "endpoint-uri", "",
		"", "Endpoint URI for the API")
	CreateDeployCmd.Flags().StringVarP(&apiSpecRevision, "external-channel-uri", "",
		"", "External Channle URI")
	CreateDeployCmd.Flags().StringVarP(&apiSpecRevision, "intended-audience", "",
		"", "Intended Audience for the API")
	CreateDeployCmd.Flags().StringVarP(&accessGuidance, "access-guidance", "",
		"", "Guidance for accessing the API")
	CreateDeployCmd.Flags().StringToStringVar(&labels, "labels",
		nil, "Labels attach identifying metadata to resources")
	CreateDeployCmd.Flags().StringToStringVar(&annotations, "annotations",
		nil, "Annotations attach non-identifying metadata to resources")

	_ = CreateDeployCmd.MarkFlagRequired("id")
	_ = CreateDeployCmd.MarkFlagRequired("name")
}

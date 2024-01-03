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

package versions

import (
	"internal/apiclient"
	"internal/client/registry/apis"

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// CreateVersionCmd to create a new instance
var CreateVersionCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a version for an API",
	Long:  "Create a version for an API Apigee Registry Instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.CreateVersion(apiVersionID, name, displayName,
			description, state, labels, annotations, primarySpec)
		return err
	},
}

var (
	apiVersionID, description, displayName, state, primarySpec string
	labels, annotations                                        map[string]string
)

func init() {
	CreateVersionCmd.Flags().StringVarP(&apiVersionID, "api-version", "",
		"", "The ID to use for the version")
	CreateVersionCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the version")
	CreateVersionCmd.Flags().StringVarP(&displayName, "display-name", "",
		"", "Human friendly display name of the version")
	CreateVersionCmd.Flags().StringVarP(&description, "desc", "d",
		"", "A detailed description")
	CreateVersionCmd.Flags().StringVarP(&state, "state", "s",
		"", "A user-definable description of the lifecycle phase of this API version")
	CreateVersionCmd.Flags().StringVarP(&primarySpec, "primary-spec", "",
		"", "The primary spec for this version; ex: projects/*/locations/*/apis/*/versions/*/specs/*")
	CreateVersionCmd.Flags().StringToStringVar(&labels, "labels",
		nil, "Labels attach identifying metadata to resources")
	CreateVersionCmd.Flags().StringToStringVar(&annotations, "annotations",
		nil, "Annotations attach non-identifying metadata to resources")

	_ = CreateVersionCmd.MarkFlagRequired("id")
}

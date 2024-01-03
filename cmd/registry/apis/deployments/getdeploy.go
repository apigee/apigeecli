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

// GetDeployCmd to get instance
var GetDeployCmd = &cobra.Command{
	Use:   "get",
	Short: "Get deployment details for an API",
	Long:  "Get deployment details for an API",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.GetDeployment(apiName, name)
		return
	},
}

func init() {
	GetDeployCmd.Flags().StringVarP(&apiName, "api-name", "",
		"", "API Name")
	GetDeployCmd.Flags().StringVarP(&name, "name", "n",
		"", "Deployment name")

	_ = GetDeployCmd.MarkFlagRequired("api-name")
}

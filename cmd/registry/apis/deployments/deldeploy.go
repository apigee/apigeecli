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

// DelDeployCmd to get instance
var DelDeployCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete deployment associated with an API",
	Long:  "Delete deployment associated with an API",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if revision {
			_, err = apis.DeleteDeploymentRevision(apiName, name)
			return
		}
		_, err = apis.DeleteDeployment(apiName, name, force)
		return
	},
}

var force, revision bool

func init() {
	DelDeployCmd.Flags().StringVarP(&apiName, "api-name", "",
		"", "API Name")
	DelDeployCmd.Flags().StringVarP(&name, "name", "n",
		"", "Deployment name")
	DelDeployCmd.Flags().BoolVarP(&force, "force", "f",
		false, "Deployment name")
	DelDeployCmd.Flags().BoolVarP(&revision, "revision", "v",
		false, "Delete a revision")

	_ = DelDeployCmd.MarkFlagRequired("api-name")
	_ = DelDeployCmd.MarkFlagRequired("name")
}

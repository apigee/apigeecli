// Copyright 2024 Google LLC
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
	"internal/client/hub"

	"github.com/spf13/cobra"
)

// UpdateCmd
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an API Deployment in API Hub",
	Long:  "Updates an API Deployment in API Hub",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = hub.UpdateDeployment(deploymentName, displayName, description, externalURI, resourceURI, endpoints, d, e, s)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&deploymentName, "id", "i",
		"", "Deployment Name")
	UpdateCmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "Deployment Display Name")
	UpdateCmd.Flags().StringVarP(&description, "description", "",
		"", "Deployment Description")
	UpdateCmd.Flags().StringVarP(&externalURI, "external-uri", "",
		"", "The uri of the externally hosted documentation")
	UpdateCmd.Flags().StringVarP(&resourceURI, "resource-uri", "",
		"", "A URI to the runtime resource")
	UpdateCmd.Flags().StringArrayVarP(&endpoints, "endpoints", "",
		[]string{}, " The endpoints at which this deployment resource is listening for API requests")
	UpdateCmd.Flags().Var(&d, "dep-type", "The type of deployment")
	UpdateCmd.Flags().Var(&e, "env-type", "The environment mapping to this deployment")
	UpdateCmd.Flags().Var(&s, "slo-type", "The SLO for this deployment")
}

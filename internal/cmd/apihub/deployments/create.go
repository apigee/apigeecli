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

// CrtCmd
var CrtCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API Deployment in API Hub",
	Long:  "Create a new API Deployment in API Hub",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = hub.CreateDeployment(deploymentID, displayName, description,
			deploymentName, externalURI, resourceURI, endpoints, d, e, s)
		return
	},
}

var (
	deploymentID, displayName, deploymentName, description, externalURI, resourceURI string
	endpoints                                                                        []string
	d                                                                                hub.DeploymentType
	e                                                                                hub.EnvironmentType
	s                                                                                hub.SloType
)

func init() {
	CrtCmd.Flags().StringVarP(&deploymentID, "id", "i",
		"", "Deployment ID")
	CrtCmd.Flags().StringVarP(&deploymentName, "name", "n",
		"", "Deployment Name")
	CrtCmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "Deployment Display Name")
	CrtCmd.Flags().StringVarP(&description, "description", "",
		"", "Deployment Description")
	CrtCmd.Flags().StringVarP(&externalURI, "external-uri", "",
		"", "The uri of the externally hosted documentation")
	CrtCmd.Flags().StringVarP(&resourceURI, "resource-uri", "",
		"", "A URI to the runtime resource")
	CrtCmd.Flags().StringArrayVarP(&endpoints, "endpoints", "",
		[]string{}, " The endpoints at which this deployment resource is listening for API requests")
	CrtCmd.Flags().Var(&d, "dep-type", "The type of deployment")
	CrtCmd.Flags().Var(&e, "env-type", "The environment mapping to this deployment")
	CrtCmd.Flags().Var(&s, "slo-type", "The SLO for this deployment")

	_ = CrtCmd.MarkFlagRequired("name")
	_ = CrtCmd.MarkFlagRequired("display-name")
	_ = CrtCmd.MarkFlagRequired("resource-uri")
	_ = CrtCmd.MarkFlagRequired("endpoints")
	_ = CrtCmd.MarkFlagRequired("type")
}

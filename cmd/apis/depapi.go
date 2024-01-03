// Copyright 2020 Google LLC
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

	"internal/client/apis"

	"github.com/spf13/cobra"
)

// DepCmd to deploy api
var DepCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a revision of an existing API proxy",
	Long: "Deploys a revision of an existing API proxy to an environment " +
		"in an organization, optionally waits for deployment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if revision == -1 {
			if revision, err = apis.GetHighestProxyRevision(name); err != nil {
				return err
			}
		}
		if _, err = apis.DeployProxy(name,
			revision,
			overrides,
			sequencedRollout,
			safeDeploy,
			serviceAccountName); err != nil {
			return err
		}

		if wait {
			err = Wait(name, revision)
		}
		return err
	},
}

var (
	overrides, wait, sequencedRollout, safeDeploy bool
	serviceAccountName                            string
)

const interval = 10

func init() {
	DepCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	DepCmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	DepCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision. If not set, the highest revision is used")
	DepCmd.Flags().BoolVarP(&overrides, "ovr", "",
		false, "Forces deployment of the new revision")
	DepCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the deployment to finish, with success or error")
	DepCmd.Flags().BoolVarP(&sequencedRollout, "sequencedrollout", "",
		false, "If set to true, the routing rules will be rolled out in a safe order; default is false")
	DepCmd.Flags().BoolVarP(&safeDeploy, "safedeploy", "",
		true, "When set to true, generateDeployChangeReport will be executed and "+
			"deployment will proceed if there are no conflicts; default is true")
	DepCmd.Flags().StringVarP(&serviceAccountName, "sa", "s",
		"", "The format must be {ACCOUNT_ID}@{PROJECT}.iam.gserviceaccount.com.")

	_ = DepCmd.MarkFlagRequired("env")
	_ = DepCmd.MarkFlagRequired("name")
}

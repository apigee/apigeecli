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
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/apis"
)

//DepCmd to deploy api
var DepCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a revision of an existing API proxy",
	Long:  "Deploys a revision of an existing API proxy to an environment in an organization",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.DeployProxy(name, revision, overrides, serviceAccountName)
		return
	},
}

var overrides bool
var serviceAccountName string

func init() {

	DepCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	DepCmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	DepCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision")
	DepCmd.Flags().BoolVarP(&overrides, "ovr", "r",
		false, "Forces deployment of the new revision")
	DepCmd.Flags().StringVarP(&serviceAccountName, "sa", "s",
		false, "The format must be {ACCOUNT_ID}@{PROJECT}.iam.gserviceaccount.com.")

	_ = DepCmd.MarkFlagRequired("env")
	_ = DepCmd.MarkFlagRequired("name")
}

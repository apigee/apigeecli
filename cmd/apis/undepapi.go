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

// UndepCmd to undeloy api
var UndepCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploys a revision of an existing API proxy",
	Long:  "Undeploys a revision of an existing API proxy to an environment in an organization",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if revision == -1 {
			if revision, err = apis.GetHighestProxyRevision(name); err != nil {
				return
			}
		}
		_, err = apis.UndeployProxy(name, revision, safeUndeploy)
		return
	},
}

var safeUndeploy bool

func init() {
	UndepCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	UndepCmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	UndepCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision")
	UndepCmd.Flags().BoolVarP(&safeUndeploy, "safeundeploy", "",
		true, "When set to true, generateUndeployChangeReport will be executed and "+
			"undeployment will proceed if there are no conflicts; default is true")

	_ = UndepCmd.MarkFlagRequired("env")
	_ = UndepCmd.MarkFlagRequired("name")
	_ = UndepCmd.MarkFlagRequired("rev")
}

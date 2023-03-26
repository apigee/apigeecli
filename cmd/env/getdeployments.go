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

package env

import (
	"internal/apiclient"

	environments "github.com/apigee/apigeecli/client/env"
	"github.com/spf13/cobra"
)

// DeployCmd to get deployed apis in an env
var GetDeployCmd = &cobra.Command{
	Use:   "get",
	Short: "Get deployments for an Environment",
	Long:  "Get deployments for an Environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if all {
			_, err = environments.GetAllDeployments()
		} else {
			_, err = environments.GetDeployments(sharedflows)
		}
		return
	},
}

var sharedflows, all bool

func init() {

	GetDeployCmd.Flags().BoolVarP(&sharedflows, "sharedflows", "s",
		false, "Return sharedflow deployments")
	GetDeployCmd.Flags().BoolVarP(&all, "all", "",
		false, "Return all deployments")
}

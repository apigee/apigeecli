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
	"fmt"

	"internal/apiclient"

	"internal/client/apis"

	"github.com/spf13/cobra"
)

// ListDepCmd to list deployed api
var ListDepCmd = &cobra.Command{
	Use:   "listdeploy",
	Short: "Lists all deployments of an API proxy",
	Long:  "Lists all deployments of an API proxy",
	Args: func(cmd *cobra.Command, args []string) error {

		apiclient.SetApigeeEnv(env)
		if apiclient.GetApigeeEnv() == "" && name == "" {
			return fmt.Errorf("proxy name or environment must be supplied")
		}
		if revision != -1 && name == "" {
			return fmt.Errorf("proxy name must be supplied with revision")
		}
		if name != "" && revision == -1 && apiclient.GetApigeeEnv() != "" {
			return fmt.Errorf("revision must be supplied with proxy name and env")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if apiclient.GetApigeeEnv() != "" {
			if revision != -1 {
				if !report {
					_, err = apis.ListProxyRevisionDeployments(name, revision)
				} else {
					_, err = apis.GenerateDeployChangeReport(name, revision, false)
				}
			} else {
				_, err = apis.ListEnvDeployments()
			}
		} else {
			_, err = apis.ListProxyDeployments(name)
		}
		return
	},
}

var report bool

func init() {
	ListDepCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	ListDepCmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	ListDepCmd.Flags().IntVarP(&revision, "rev", "r",
		-1, "API Proxy revision")
	ListDepCmd.Flags().BoolVarP(&report, "report", "",
		false, "Generate Deploy Change Report; Default is false")
}

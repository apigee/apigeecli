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

package sharedflows

import (
	"fmt"
	"internal/apiclient"
	"internal/client/sharedflows"

	"github.com/spf13/cobra"
)

// ListDepCmd to list deployed api
var ListDepCmd = &cobra.Command{
	Use:   "listdeploy",
	Short: "Lists all deployments of a Sharedflow",
	Long:  "Lists all deployments of a Sharedflow",
	Args: func(cmd *cobra.Command, args []string) error {
		apiclient.SetApigeeEnv(env)
		if apiclient.GetApigeeEnv() == "" && name == "" {
			return fmt.Errorf("sharedflow name or environment must be supplied")
		}
		if revision != -1 && name == "" {
			return fmt.Errorf("sharedflow name must be supplied with revision")
		}
		if name != "" && revision == -1 && apiclient.GetApigeeEnv() != "" {
			return fmt.Errorf("revision must be supplied with sharedflow name and env")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if apiclient.GetApigeeEnv() != "" {
			if revision != -1 {
				_, err = sharedflows.ListRevisionDeployments(name, revision)
			} else {
				_, err = sharedflows.ListEnvDeployments()
			}
		} else {
			_, err = sharedflows.ListDeployments(name)
		}
		return err
	},
}

func init() {
	ListDepCmd.Flags().StringVarP(&name, "name", "n",
		"", "Shareflow name")
	ListDepCmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	ListDepCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "Shareflow revision")
}

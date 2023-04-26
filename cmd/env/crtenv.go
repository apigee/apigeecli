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

	"internal/client/env"

	"github.com/spf13/cobra"
)

// CreateCmd to create env
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Long:  "Create a new environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = env.Create(deploymentType, apiProxyType)
		return
	},
}

var deploymentType, apiProxyType string

func init() {
	CreateCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")
	CreateCmd.Flags().StringVarP(&deploymentType, "deptype", "d",
		"", "Deployment type - must be PROXY or ARCHIVE")
	CreateCmd.Flags().StringVarP(&apiProxyType, "proxtype", "p",
		"", "Proxy type - must be PROGRAMMABLE or CONFIGURABLE")
	_ = CreateCmd.MarkFlagRequired("env")
}

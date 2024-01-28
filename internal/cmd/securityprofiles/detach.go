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

package securityprofiles

import (
	"internal/apiclient"
	"internal/client/securityprofiles"

	"github.com/spf13/cobra"
)

// DetachCmd to list catalog items
var DetachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach a security profile from an environment",
	Long:  "Detach a security profile from an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = securityprofiles.Detach(name)
		return
	},
}

func init() {
	DetachCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the security profile")
	DetachCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")
	_ = DetachCmd.MarkFlagRequired("name")
	_ = DetachCmd.MarkFlagRequired("env")
}

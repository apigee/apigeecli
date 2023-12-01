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

// AttachCmd to list catalog items
var AttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a security profile to an environment",
	Long:  "Attach a security profile to an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = securityprofiles.Attach(name, revision)
		return
	},
}

var environment, revision string

func init() {
	AttachCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the security profile")
	AttachCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")
	AttachCmd.Flags().StringVarP(&revision, "rev", "r",
		"", "Security Profile revision id")
	_ = AttachCmd.MarkFlagRequired("name")
	_ = AttachCmd.MarkFlagRequired("env")
	_ = AttachCmd.MarkFlagRequired("rev")
}

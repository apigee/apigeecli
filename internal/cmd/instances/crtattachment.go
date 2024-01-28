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

package instances

import (
	"internal/apiclient"

	"internal/client/instances"

	"github.com/spf13/cobra"
)

// CreateAttachCmd to create a new instance
var CreateAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach an environment to an instance",
	Long:  "Attach an environment to an instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = instances.Attach(name, environment)
		return
	},
}

func init() {
	CreateAttachCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Instance")
	CreateAttachCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")

	_ = CreateAttachCmd.MarkFlagRequired("name")
	_ = CreateAttachCmd.MarkFlagRequired("env")
}

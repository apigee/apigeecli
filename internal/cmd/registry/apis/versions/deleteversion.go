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

package versions

import (
	"internal/apiclient"

	"internal/client/registry/apis"

	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// DelVersionCmd to get instance
var DelVersionCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an API Version for an API",
	Long:  "Delete an API Version for an API in Apigee Registry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.DeleteVersion(name, version, force)
		return
	},
}

var force bool

func init() {
	DelVersionCmd.Flags().StringVarP(&name, "name", "n",
		"", "Apigee Registry API Name")
	DelVersionCmd.Flags().StringVarP(&version, "version", "v",
		"", "API Version name")
	DelVersionCmd.Flags().BoolVarP(&force, "force", "f",
		false, "Deployment name")

	_ = DelVersionCmd.MarkFlagRequired("name")
	_ = DelVersionCmd.MarkFlagRequired("version")
}

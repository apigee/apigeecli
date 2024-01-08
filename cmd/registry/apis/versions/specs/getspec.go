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

package specs

import (
	"internal/apiclient"

	"internal/client/registry/apis/versions"

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// GetSpecCmd to get instance
var GetSpecCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the spec details for an API version",
	Long:  "Get the spec details for an API version in Apigee Registry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if content {
			_, err = versions.GetSpecContents(apiName, apiVersion, name)
			return err
		}
		_, err = versions.GetSpec(apiName, apiVersion, name)
		return
	},
}

var (
	name    string
	content bool
)

func init() {
	GetSpecCmd.Flags().StringVarP(&apiName, "api-name", "",
		"", "API Name")
	GetSpecCmd.Flags().StringVarP(&apiVersion, "api-version", "",
		"", "API Version")
	GetSpecCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Spec")
	GetSpecCmd.Flags().BoolVarP(&content, "content", "c",
		false, "If set to true, returns artifact contents")

	_ = GetSpecCmd.MarkFlagRequired("api-name")
	_ = GetSpecCmd.MarkFlagRequired("api-version")
	_ = GetSpecCmd.MarkFlagRequired("name")
}

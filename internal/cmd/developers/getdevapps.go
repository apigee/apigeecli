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

package developers

import (
	"internal/apiclient"

	"internal/client/developers"

	"github.com/spf13/cobra"
)

// GetAppCmd to get developer aps
var GetAppCmd = &cobra.Command{
	Use:   "getapps",
	Short: "Returns the apps owned by a developer by email address",
	Long:  "Returns the apps owned by a developer by email address",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = developers.GetApps(name, expand)
		return
	},
}

var name string

func init() {
	GetAppCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	GetAppCmd.Flags().StringVarP(&name, "name", "n",
		"", "email of the developer")

	GetAppCmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "expand app details")

	_ = GetAppCmd.MarkFlagRequired("org")
	_ = GetAppCmd.MarkFlagRequired("name")
}

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

package res

import (
	"internal/apiclient"
	"internal/client/res"
	"strings"

	"github.com/spf13/cobra"
)

// UpdateCmd to get a resource
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a resource file",
	Long:  "Update a resource file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		return res.Update(name, resPath, resType)
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the resource file")
	UpdateCmd.Flags().StringVarP(&resType, "type", "p",
		"", "Resource type. Valid types include "+strings.Join(res.GetValidResourceTypes(), ", "))
	UpdateCmd.Flags().StringVarP(&resPath, "respath", "",
		"", "Resource Path")

	_ = UpdateCmd.MarkFlagRequired("name")
	_ = UpdateCmd.MarkFlagRequired("type")
	_ = UpdateCmd.MarkFlagRequired("respath")
}

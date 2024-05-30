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
	"strings"

	"internal/client/res"

	"github.com/spf13/cobra"
)

// CreateCmd to create a resource
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource file",
	Long:  "Create a resource file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = res.Create(name, resPath, resType)
		return
	},
}

var resPath string

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the resource file")
	CreateCmd.Flags().StringVarP(&resType, "type", "p",
		"", "Resource type; Valid types include "+strings.Join(res.GetValidResourceTypes(), ", "))
	CreateCmd.Flags().StringVarP(&resPath, "respath", "",
		"", "Resource Path")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("type")
	_ = CreateCmd.MarkFlagRequired("respath")
}

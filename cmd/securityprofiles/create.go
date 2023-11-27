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

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// CreateCmd to get a securityprofile
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Security Profile",
	Long:  "Create a new Security Profile",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		content, err := utils.ReadFile(securityActionFile)
		if err != nil {
			return err
		}
		_, err = securityprofiles.Create(name, content)
		return
	},
}

var securityActionFile string

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Security Action name")
	CreateCmd.Flags().StringVarP(&securityActionFile, "file", "f",
		"", "Path to a file containing Security Profile content")
	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("file")
}

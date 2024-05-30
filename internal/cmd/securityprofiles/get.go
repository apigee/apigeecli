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

// GetCmd to list catalog items
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns a security profile by name",
	Long:  "Returns a security profile by name",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = securityprofiles.Get(name, revision)
		return
	},
}

var name string

func init() {
	GetCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the security profile")
	GetCmd.Flags().StringVarP(&revision, "revision", "v",
		"", "Revision of the security profile")
	_ = GetCmd.MarkFlagRequired("name")
}

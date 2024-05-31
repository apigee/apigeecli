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

package envgroup

import (
	"internal/apiclient"

	"internal/client/envgroups"

	"github.com/spf13/cobra"
)

// UpdateCmd to create a new product
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update hostnames in an Environment Group",
	Long:  "Update hostnames in an Environment Group",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = envgroups.PatchHosts(name, hostnames)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Environment Group")
	UpdateCmd.Flags().StringArrayVarP(&hostnames, "hosts", "d",
		[]string{}, "A list of hostnames")

	_ = UpdateCmd.MarkFlagRequired("name")
	_ = UpdateCmd.MarkFlagRequired("hosts")
}

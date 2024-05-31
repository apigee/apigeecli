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

// UpdateCmd to create a new instance
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an Instance",
	Long:  "Update an Instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = instances.Update(name, consumerAcceptList)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Instance")
	UpdateCmd.Flags().StringArrayVarP(&consumerAcceptList, "consumer-accept-list", "c",
		[]string{}, "Customer accept list represents the list of "+
			"projects (id/number) that can connect to the service attachment")

	_ = UpdateCmd.MarkFlagRequired("name")
	_ = UpdateCmd.MarkFlagRequired("consumer-accept-list")
}

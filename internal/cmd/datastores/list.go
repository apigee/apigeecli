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

package datastores

import (
	"internal/apiclient"
	"internal/client/datastores"

	"github.com/spf13/cobra"
)

// ListCmd to list datastores
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of data stores in the org",
	Long:  "Returns a list of data stores in the org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = datastores.List(targetType)
		return
	},
}

func init() {
	ListCmd.Flags().StringVarP(&targetType, "targettype", "",
		"", "TargetType is used to fetch datastores of matching type; default is fetch all")
}

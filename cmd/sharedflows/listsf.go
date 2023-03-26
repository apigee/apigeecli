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

package sharedflows

import (
	"internal/apiclient"

	"github.com/apigee/apigeecli/client/sharedflows"
	"github.com/spf13/cobra"
)

// Cmd to list shared flow
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all shared flows in the organization.",
	Long:  "Lists all shared flows in the organization.",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = sharedflows.List(includeRevisions)
		return
	},
}

var includeRevisions bool

func init() {

	ListCmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	ListCmd.Flags().BoolVarP(&includeRevisions, "rev", "r",
		false, "Include revisions")

}

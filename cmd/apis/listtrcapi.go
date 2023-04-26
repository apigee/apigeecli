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

package apis

import (
	"internal/apiclient"

	"internal/client/apis"

	"github.com/spf13/cobra"
)

// ListTrcCmd to manage tracing of apis
var ListTrcCmd = &cobra.Command{
	Use:   "list",
	Short: "List all debug sessions for an API proxy revision",
	Long:  "List all debug sessions for an API proxy revision deployed in an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.ListTracceSession(name, revision)
		return
	},
}

func init() {
	ListTrcCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	ListTrcCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision")

	_ = ListTrcCmd.MarkFlagRequired("name")
	_ = ListTrcCmd.MarkFlagRequired("rev")
}

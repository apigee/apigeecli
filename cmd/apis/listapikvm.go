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

// ListKvmCmd to manage tracing of apis
var ListKvmCmd = &cobra.Command{
	Use:   "list",
	Short: "List all KVMs for an API proxy",
	Long:  "List all KVMs for an API proxy",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.ListProxyKVM(proxyName)
		return
	},
}

func init() {
	ListKvmCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")

	_ = ListKvmCmd.MarkFlagRequired("proxy")
}

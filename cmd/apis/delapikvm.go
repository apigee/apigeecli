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
	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/apis"
	"github.com/spf13/cobra"
)

//DelKvmCmd to manage tracing of apis
var DelKvmCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an API Proxy scoped KVM",
	Long:  "Deletes an API Proxy scoped KVM",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.DeleteProxyKVM(proxyName, name)
		return
	},
}

func init() {
	DelKvmCmd.Flags().StringVarP(&name, "name", "n",
		"", "KVM name")
	DelKvmCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")

	_ = DelKvmCmd.MarkFlagRequired("proxy")
	_ = DelKvmCmd.MarkFlagRequired("name")
}

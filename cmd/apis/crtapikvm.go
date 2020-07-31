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
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/apis"
)

//CreateKvmCmd to manage tracing of apis
var CreateKvmCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an API Proxy scoped KVM",
	Long:  "Create an API Proxy scoped KVM",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.CreateProxyKVM(proxyName, name, encrypted)
		return
	},
}

var encrypted bool

func init() {
	CreateKvmCmd.Flags().StringVarP(&name, "name", "n",
		"", "KVM name")
	CreateKvmCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")
	CreateKvmCmd.Flags().BoolVarP(&encrypted, "env", "e",
		true, "Encrypted")

	_ = CreateKvmCmd.MarkFlagRequired("proxy")
	_ = CreateKvmCmd.MarkFlagRequired("name")
}

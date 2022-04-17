// Copyright 2021 Google LLC
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

package kvm

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/kvm"
)

//DelEntryCmd to delete kvm map entry
var DelEntryCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a KVM Map entry",
	Long:  "Delete a KVM Map entry",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if env != "" {
			apiclient.SetApigeeEnv(env)
		}
		if env != "" && proxyName != "" {
			return fmt.Errorf("proxy and env flags cannot be used together")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = kvm.GetEntry(proxyName, mapName, keyName)
		return
	},
}

func init() {
	DelEntryCmd.Flags().StringVarP(&mapName, "map", "m",
		"", "KV Map Name")
	DelEntryCmd.Flags().StringVarP(&env, "env", "e",
		"", "Environment name")
	DelEntryCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")
	DelEntryCmd.Flags().StringVarP(&keyName, "key", "k",
		"", "KV Map entry name")

	_ = DelEntryCmd.MarkFlagRequired("key")
	_ = DelEntryCmd.MarkFlagRequired("map")
}

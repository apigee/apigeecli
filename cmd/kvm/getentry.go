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

//GetEntryCmd to create kvm map entry
var GetEntryCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a KV Map entry",
	Long:  "Get a KV Map entry",
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
	GetEntryCmd.Flags().StringVarP(&mapName, "map", "m",
		"", "KV Map Name")
	GetEntryCmd.Flags().StringVarP(&env, "env", "e",
		"", "Environment name")
	GetEntryCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")
	GetEntryCmd.Flags().StringVarP(&keyName, "key", "k",
		"", "KV Map entry name")

	_ = GetEntryCmd.MarkFlagRequired("key")
	_ = GetEntryCmd.MarkFlagRequired("map")
}

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

//CreateEntryCmd to create kv map entry
var CreateEntryCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a KV Map entry",
	Long:  "Create a KV Map entry",
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
		_, err = kvm.CreateEntry(proxyName, mapName, keyName, value)
		return
	},
}

var value string

func init() {
	CreateEntryCmd.Flags().StringVarP(&mapName, "map", "m",
		"", "KV Map Name")
	CreateEntryCmd.Flags().StringVarP(&env, "env", "e",
		"", "Environment name")
	CreateEntryCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")
	CreateEntryCmd.Flags().StringVarP(&keyName, "key", "k",
		"", "KV Map entry name")
	CreateEntryCmd.Flags().StringVarP(&value, "value", "l",
		"", "KV Map entry value")

	_ = CreateEntryCmd.MarkFlagRequired("key")
	_ = CreateEntryCmd.MarkFlagRequired("value")
	_ = CreateEntryCmd.MarkFlagRequired("map")
}

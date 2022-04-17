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

//ListEntryCmd to list kvm map entries
var ListEntryCmd = &cobra.Command{
	Use:   "list",
	Short: "List KV Map entries",
	Long:  "List KV Map entries",
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
		_, err = kvm.ListEntries(proxyName, mapName, pageSize, pageToken)
		return
	},
}

var pageToken string
var pageSize int

func init() {
	ListEntryCmd.Flags().StringVarP(&mapName, "map", "m",
		"", "KV Map Name")
	ListEntryCmd.Flags().StringVarP(&env, "env", "e",
		"", "Environment name")
	ListEntryCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")
	ListEntryCmd.Flags().StringVarP(&pageToken, "page-token", "",
		"", "next_page_token from the prior response to be used to fetch the next dataset")
	ListEntryCmd.Flags().IntVarP(&pageSize, "page-size", "",
		-1, "Number of items to return on the list")

	_ = ListEntryCmd.MarkFlagRequired("map")
}

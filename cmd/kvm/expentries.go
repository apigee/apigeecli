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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/kvm"
)

//ExpEntryCmd to export map entries to files
var ExpEntryCmd = &cobra.Command{
	Use:   "export",
	Short: "Export KV Map entries",
	Long:  "Export KV Map entries",
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
		var payload [][]byte
		var fileName string

		if payload, err = kvm.ExportEntries(proxyName, mapName); err != nil {
			return
		}

		if env != "" {
			fileName = strings.Join([]string{"env", env, mapName, "kvmfile"}, "_")
		} else if proxyName != "" {
			fileName = strings.Join([]string{"proxy", proxyName, mapName, "kvmfile"}, "_")
		} else {
			fileName = strings.Join([]string{"org_", mapName, "kvmfile"}, "_")
		}

		for i := range payload {
			if err = apiclient.WriteByteArrayToFile(fileName+"_"+strconv.Itoa(i)+".json", false, payload[i]); err != nil {
				return
			}
		}
		return
	},
}

func init() {
	ExpEntryCmd.Flags().StringVarP(&mapName, "map", "m",
		"", "KV Map Name")
	ExpEntryCmd.Flags().StringVarP(&env, "env", "e",
		"", "Environment name")
	ExpEntryCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")

	_ = ExpEntryCmd.MarkFlagRequired("map")
}

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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"internal/cmd/utils"

	"internal/apiclient"
	"internal/clilog"

	"internal/client/kvm"

	"github.com/spf13/cobra"
)

// ExpCmd to export map entries to files
var ExpCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all KV Map entries for all KV Maps",
	Long:  "Export all KV Map entries for all KV Maps in a given scope",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if env != "" {
			apiclient.SetApigeeEnv(env)
		}
		if env != "" && proxyName != "" {
			return fmt.Errorf("proxy and env flags cannot be used together")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var payload [][]byte
		var fileName string

		apiclient.DisableCmdPrintHttpResponse()

		// return all kvm entries from all proxies
		if env == "" && proxyName == "" {
			return kvm.ExportAllEntries()
		}

		listKVMBytes, err := kvm.List(proxyName)
		if err != nil {
			return err
		}

		var listKVM []string
		if err = json.Unmarshal(listKVMBytes, &listKVM); err != nil {
			return err
		}

		for _, mapName := range listKVM {
			if payload, err = kvm.ExportEntries(proxyName, mapName); err != nil {
				return err
			}

			if env != "" {
				fileName = strings.Join([]string{"env", env, mapName, "kvmfile"}, utils.DefaultFileSplitter)
			} else if proxyName != "" {
				fileName = strings.Join([]string{"proxy", proxyName, mapName, "kvmfile"}, utils.DefaultFileSplitter)
			} else {
				fileName = strings.Join([]string{"org", mapName, "kvmfile"}, utils.DefaultFileSplitter)
			}

			for i := range payload {
				if err = apiclient.WriteByteArrayToFile(fileName+utils.DefaultFileSplitter+strconv.Itoa(i)+".json", false, payload[i]); err != nil {
					return err
				}
			}
		}

		clilog.Info.Println("KVMs exported successfully")
		return err
	},
}

func init() {
	ExpCmd.Flags().StringVarP(&env, "env", "e",
		"", "Environment name")
	ExpCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")
}

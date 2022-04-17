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

//ImpEntryCmd to import kvm entries from files
var ImpEntryCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing KVM Entries",
	Long:  "Import a file containing KVM Entries",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if env != "" {
			apiclient.SetApigeeEnv(env)
		}
		if env != "" && proxyName != "" {
			return fmt.Errorf("proxy and env flags cannot be used together")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return kvm.ImportEntries(proxyName, mapName, conn, filePath)
	},
}

var conn int
var filePath string

func init() {

	ImpEntryCmd.Flags().StringVarP(&filePath, "file", "f",
		"", "File containing App Developers")
	ImpEntryCmd.Flags().StringVarP(&mapName, "map", "m",
		"", "KV Map Name")
	ImpEntryCmd.Flags().StringVarP(&env, "env", "e",
		"", "Environment name")
	ImpEntryCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")
	ImpEntryCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = ImpEntryCmd.MarkFlagRequired("map")
	_ = ImpEntryCmd.MarkFlagRequired("file")
}

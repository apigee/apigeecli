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
	"strings"

	"internal/apiclient"
	"internal/clilog"

	"internal/client/kvm"

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// ImpCmd to import kvm entries from files
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import KVM Entries from a folder containing KVM files",
	Long:  "Import KVM Entries from a folder containing KVM files",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.DisableCmdPrintHttpResponse()

		orgKVMFileList, envKVMFileList, proxyKVMFileList, err := utils.ListKVMFiles(folder, skipExistingKVMs)
		if err != nil {
			return err
		}

		if len(orgKVMFileList) > 0 {
			clilog.Info.Println("Importing org scoped KVMs...")
			for _, orgKVMFile := range orgKVMFileList {
				kvmMetadata := strings.Split(orgKVMFile, "_")
				clilog.Info.Printf("\tCreating KVM %s\n", orgKVMFile)
				if _, err = kvm.Create("", kvmMetadata[1], true); err != nil {
					return err
				}
				clilog.Info.Printf("\tImporting entries for %s\n", orgKVMFile)
				if err = kvm.ImportEntries("", kvmMetadata[1], conn, orgKVMFile); err != nil {
					return err
				}
			}
		}

		if len(envKVMFileList) > 0 {
			clilog.Info.Println("Importing env scoped KVMs...")
			for _, envKVMFile := range envKVMFileList {
				kvmMetadata := strings.Split(envKVMFile, "_")
				apiclient.SetApigeeEnv(kvmMetadata[1])
				clilog.Info.Printf("\tCreating KVM %s\n", envKVMFile)
				if _, err = kvm.Create("", kvmMetadata[2], true); err != nil {
					return err
				}
				clilog.Info.Printf("\tImporting entries for %s\n", envKVMFile)
				if err = kvm.ImportEntries("", kvmMetadata[2], conn, envKVMFile); err != nil {
					return err
				}
			}
		}

		if len(proxyKVMFileList) > 0 {
			clilog.Info.Println("Importing proxy scoped KVMs...")
			for _, proxyKVMFile := range proxyKVMFileList {
				kvmMetadata := strings.Split(proxyKVMFile, "_")
				clilog.Info.Printf("\tCreating KVM %s\n", proxyKVMFile)
				if _, err = kvm.Create(kvmMetadata[1], kvmMetadata[2], true); err != nil {
					return err
				}
				clilog.Info.Printf("\tImporting entries for %s\n", proxyKVMFile)
				if err = kvm.ImportEntries(kvmMetadata[1], kvmMetadata[2], conn, proxyKVMFile); err != nil {
					return err
				}
			}
		}

		return err
	},
}

var (
	folder           string
	skipExistingKVMs bool
)

func init() {
	ImpCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "The absolute path to the folder containing KVM entries")
	ImpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
	ImpCmd.Flags().BoolVarP(&skipExistingKVMs, "skipExistingKVMs", "",
		false, "Skip import of existing KVM(s); default false")

	_ = ImpCmd.MarkFlagRequired("folder")
}

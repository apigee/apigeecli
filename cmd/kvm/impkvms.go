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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/kvm"
	"github.com/spf13/cobra"
)

// ImpCmd to import kvm entries from files
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing KVM Entries",
	Long:  "Import a file containing KVM Entries",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		orgKVMFileList, envKVMFileList, proxyKVMFileList, err := listKVMFiles()
		if err != nil {
			return err
		}
		if len(orgKVMFileList) > 0 {
			fmt.Println("Importing org scoped KVMs...")
			for _, orgKVMFile := range orgKVMFileList {
				fmt.Printf("\tImporting %s\n", orgKVMFile)
				kvmMetadata := strings.Split(orgKVMFile, "_")
				if err = kvm.ImportEntries("", kvmMetadata[1], conn, orgKVMFile); err != nil {
					return err
				}
			}
		}

		if len(envKVMFileList) > 0 {
			fmt.Println("Importing env scoped KVMs...")
			for _, envKVMFile := range envKVMFileList {
				fmt.Printf("\tImporting %s\n", envKVMFile)
				kvmMetadata := strings.Split(envKVMFile, "_")
				apiclient.SetApigeeEnv(kvmMetadata[1])
				if err = kvm.ImportEntries("", kvmMetadata[2], conn, envKVMFile); err != nil {
					return err
				}
			}
		}

		if len(proxyKVMFileList) > 0 {
			fmt.Println("Importing proxy scoped KVMs...")
			for _, proxyKVMFile := range proxyKVMFileList {
				fmt.Printf("\tImporting %s\n", proxyKVMFile)
				kvmMetadata := strings.Split(proxyKVMFile, "_")
				if err = kvm.ImportEntries(kvmMetadata[1], kvmMetadata[2], conn, proxyKVMFile); err != nil {
					return err
				}
			}
		}

		return err
	},
}

var folder string

func init() {

	ImpCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "The absolute path to the folder containing KVM entries")
	ImpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = ImpCmd.MarkFlagRequired("folder")
}

func listKVMFiles() (orgKVMFileList map[string]string, envKVMFileList map[string]string, proxyKVMFileList map[string]string, err error) {

	orgKVMFileList = map[string]string{}
	envKVMFileList = map[string]string{}
	proxyKVMFileList = map[string]string{}

	renv := regexp.MustCompile(`env_(\S*)_kvmfile_[0-9]+\.json`)
	rorg := regexp.MustCompile(`org_(\S*)_kvmfile_[0-9]+\.json`)
	rproxy := regexp.MustCompile(`proxy_(\S*)_kvmfile_[0-9]+\.json`)

	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			kvmFile := filepath.Base(path)
			if renv.MatchString(kvmFile) {
				envKVMFileSplit := strings.Split(kvmFile, "_")
				if len(envKVMFileSplit) > 2 {
					fmt.Printf("Map name %s, path %s\n", envKVMFileSplit[2], kvmFile)
					envKVMFileList[envKVMFileSplit[2]] = path
				}
			} else if rproxy.MatchString(kvmFile) {
				proxyKVMFileSplit := strings.Split(kvmFile, "_")
				if len(proxyKVMFileSplit) > 2 {
					fmt.Printf("Map name %s, path %s\n", proxyKVMFileSplit[2], kvmFile)
					proxyKVMFileList[proxyKVMFileSplit[2]] = path
				}
			} else if rorg.MatchString(kvmFile) {
				orgKVMFileSplit := strings.Split(kvmFile, "_")
				if len(orgKVMFileSplit) > 1 {
					fmt.Printf("Map name %s, path %s\n", orgKVMFileSplit[1], kvmFile)
					orgKVMFileList[orgKVMFileSplit[1]] = path
				}
			}
		}
		return nil
	})
	return orgKVMFileList, envKVMFileList, proxyKVMFileList, err
}

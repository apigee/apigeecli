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

package org

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/apis"
	"github.com/srinandan/apigeecli/client/apps"
	"github.com/srinandan/apigeecli/client/developers"
	"github.com/srinandan/apigeecli/client/env"
	"github.com/srinandan/apigeecli/client/envgroups"
	"github.com/srinandan/apigeecli/client/keystores"
	"github.com/srinandan/apigeecli/client/kvm"
	"github.com/srinandan/apigeecli/client/products"
	"github.com/srinandan/apigeecli/client/sharedflows"
	"github.com/srinandan/apigeecli/client/targetservers"
	"github.com/srinandan/apigeecli/clilog"
)

//ImportCmd to get org details
var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Apigee Configuration",
	Long:  "Import Apigee Configuration",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		var keystoreList, kvmList []string

		clilog.Warning.Println("Calls to Apigee APIs have a quota of 6000 per min. Running this tool against large list of entities can exhaust that quota and impact the usage of the platform.")

		fmt.Println("Importing API Proxies...")
		if err = apis.ImportProxies(conn, path.Join(folder, proxiesFolderName)); err != nil {
			return err
		}

		fmt.Println("Importing Sharedflows...")
		if err = sharedflows.Import(conn, path.Join(folder, sharedFlowsFolderName)); err != nil {
			return err
		}

		fmt.Println("Check for files with KVM Entries")
		orgKVMFileList, envKVMFileList, _, _ := listKVMFiles()

		if isFileExists(path.Join(folder, org+"_"+kVMFileName)) {
			fmt.Println("Importing Org scoped KVMs...")
			if kvmList, err = readEntityFile(path.Join(folder, org+"_"+kVMFileName)); err != nil {
				return err
			}
			for _, kvmName := range kvmList {
				//create only encrypted KVMs
				if _, err = kvm.Create("", kvmName, true); err != nil {
					return err
				}
				if orgKVMFileList[kvmName] != "" {
					kvm.ImportEntries("", kvmName, conn, orgKVMFileList[kvmName])
				}
			}
		}

		if isFileExists(path.Join(folder, productsFileName)) {
			fmt.Println("Importing Products...")
			if err = products.Import(conn, path.Join(folder, productsFileName)); err != nil {
				return err
			}
		}

		if isFileExists(path.Join(folder, developersFileName)) {
			fmt.Println("Importing Developers...")
			if err = developers.Import(conn, path.Join(folder, developersFileName)); err != nil {
				return err
			}

			fmt.Println("Importing Apps...")
			if err = apps.Import(conn, path.Join(folder, appsFileName), path.Join(folder, developersFileName)); err != nil {
				return err
			}
		}

		if isFileExists(path.Join(folder, envGroupsFileName)) {
			fmt.Println("Importing Environment Group Configuration...")
			if err = envgroups.Import(path.Join(folder, envGroupsFileName)); err != nil {
				return err
			}
		}

		apiclient.SetPrintOutput(false)

		var envRespBody []byte
		if envRespBody, err = env.List(); err != nil {
			return err
		}

		environments := []string{}

		if envRespBody != nil {
			if err = json.Unmarshal(envRespBody, &environments); err != nil {
				return err
			}
		}

		for _, environment := range environments {
			fmt.Println("Importing configuration for environment " + environment)
			apiclient.SetApigeeEnv(environment)

			if isFileExists(path.Join(folder, keyStoresFileName)) {
				fmt.Println("\tImporting Key stores...")
				if keystoreList, err = readEntityFile(path.Join(folder, keyStoresFileName)); err != nil {
					return err
				}
				for _, keystore := range keystoreList {
					if _, err = keystores.Create(keystore); err != nil {
						return err
					}
				}
			}

			if isFileExists(path.Join(folder, targetServerFileName)) {
				fmt.Println("\tImporting Target servers...")
				if err = targetservers.Import(conn, path.Join(folder, targetServerFileName)); err != nil {
					return err
				}
			}

			if isFileExists(path.Join(folder, kVMFileName)) {
				fmt.Println("\tImporting KVM Names...")
				if kvmList, err = readEntityFile(path.Join(folder, environment+"_"+kVMFileName)); err != nil {
					return err
				}
				for _, kvmName := range kvmList {
					//create only encrypted KVMs
					if _, err = kvm.Create("", kvmName, true); err != nil {
						return err
					}
					if envKVMFileList[kvmName] != "" {
						kvm.ImportEntries("", kvmName, conn, envKVMFileList[kvmName])
					}
				}
			}

			if importDebugmask {
				if isFileExists(path.Join(folder, environment+debugmaskFileName)) {
					fmt.Println("\tImporting Debug Mask configuration...")
					debugMask, _ := readEntityFileAsString(path.Join(folder, environment+debugmaskFileName))
					if _, err = env.SetDebug(debugMask); err != nil {
						return err
					}
				}
			}

			if importTrace {
				if isFileExists(path.Join(folder, environment+tracecfgFileName)) {
					fmt.Println("\tImporting Trace configuration...")
					traceCfg, _ := readEntityFileAsString(path.Join(folder, environment+tracecfgFileName))
					if _, err = env.ImportTraceConfig(traceCfg); err != nil {
						return err
					}
				}
			}
		}

		return err
	},
}

var importTrace, importDebugmask bool
var folder string

func init() {

	ImportCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	ImportCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
	ImportCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "folder containing API proxy bundles")
	ImportCmd.Flags().BoolVarP(&importTrace, "importTrace", "",
		false, "Import distributed trace configuration; default false")
	ImportCmd.Flags().BoolVarP(&importDebugmask, "importDebugmask", "",
		false, "Import debugmask configuration; default false")

	_ = ImportCmd.MarkFlagRequired("folder")
}

func isFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func readEntityFile(filePath string) ([]string, error) {

	entities := []string{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return entities, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return entities, err
	}

	if err = json.Unmarshal(byteValue, &entities); err != nil {
		return entities, err
	}

	return entities, nil
}

func readEntityFileAsString(filePath string) (string, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	return string(byteValue[:]), nil
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
			if renv.MatchString(filepath.Base(path)) {
				envKVMFileSplit := strings.Split(path, "_")
				if len(envKVMFileSplit) > 2 {
					fmt.Printf("Map name %s, path %s\n", envKVMFileSplit[2], path)
					envKVMFileList[envKVMFileSplit[2]] = path
				}
			} else if rproxy.MatchString(filepath.Base(path)) {
				proxyKVMFileSplit := strings.Split(path, "_")
				if len(proxyKVMFileSplit) > 2 {
					fmt.Printf("Map name %s, path %s\n", proxyKVMFileSplit[2], path)
					proxyKVMFileList[proxyKVMFileSplit[2]] = path
				}
			} else if rorg.MatchString(filepath.Base(path)) {
				orgKVMFileSplit := strings.Split(path, "_")
				if len(orgKVMFileSplit) > 1 {
					fmt.Printf("Map name %s, path %s\n", orgKVMFileSplit[1], path)
					orgKVMFileList[orgKVMFileSplit[1]] = path
				}
			}
		}
		return nil
	})
	return orgKVMFileList, envKVMFileList, proxyKVMFileList, err
}

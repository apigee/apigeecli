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
	"io"
	"os"
	"path"

	"internal/apiclient"

	"internal/clilog"

	"internal/client/apis"
	"internal/client/apps"
	"internal/client/datacollectors"
	"internal/client/developers"
	"internal/client/env"
	"internal/client/envgroups"
	"internal/client/keystores"
	"internal/client/kvm"
	"internal/client/products"
	"internal/client/references"
	"internal/client/sharedflows"
	"internal/client/targetservers"

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// ImportCmd to get org details
var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Apigee Configuration",
	Long:  "Import Apigee Configuration",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		var kvmList []string

		if stat, err := os.Stat(folder); err == nil && !stat.IsDir() {
			return fmt.Errorf("supplied path is not a folder")
		}

		clilog.Warning.Println("Calls to Apigee APIs have a quota of 6000 per min. Running this tool against large list of entities can exhaust that quota and impact the usage of the platform.")

		clilog.Info.Println("Importing API Proxies...")
		if err = apis.ImportProxies(conn, path.Join(folder, proxiesFolderName)); err != nil {
			return err
		}

		clilog.Info.Println("Importing Sharedflows...")
		if err = sharedflows.Import(conn, path.Join(folder, sharedFlowsFolderName)); err != nil {
			return err
		}

		clilog.Info.Println("Check for files with KVM Entries")
		orgKVMFileList, envKVMFileList, _, _ := utils.ListKVMFiles(folder)

		if utils.FileExists(path.Join(folder, "org_"+org+"_"+kVMFileName)) {
			clilog.Info.Println("Importing Org scoped KVMs...")
			if kvmList, err = utils.ReadEntityFile(path.Join(folder, "org_"+org+"_"+kVMFileName)); err != nil {
				return err
			}
			for _, kvmName := range kvmList {
				//create only encrypted KVMs
				if _, err = kvm.Create("", kvmName, true); err != nil {
					return err
				}
				if orgKVMFileList[kvmName] != "" {
					if err = kvm.ImportEntries("", kvmName, conn, orgKVMFileList[kvmName]); err != nil {
						return err
					}
				}
			}
		}

		if utils.FileExists(path.Join(folder, productsFileName)) {
			clilog.Info.Println("Importing Products...")
			if err = products.Import(conn, path.Join(folder, productsFileName), false); err != nil {
				return err
			}
		}

		if utils.FileExists(path.Join(folder, developersFileName)) {
			clilog.Info.Println("Importing Developers...")
			if err = developers.Import(conn, path.Join(folder, developersFileName)); err != nil {
				return err
			}

			clilog.Info.Println("Importing Apps...")
			if err = apps.Import(conn, path.Join(folder, appsFileName), path.Join(folder, developersFileName)); err != nil {
				return err
			}
		}

		if utils.FileExists(path.Join(folder, envGroupsFileName)) {
			clilog.Info.Println("Importing Environment Group Configuration...")
			if err = envgroups.Import(path.Join(folder, envGroupsFileName)); err != nil {
				return err
			}
		}

		if utils.FileExists(path.Join(folder, dataCollFileName)) {
			clilog.Info.Println("Importing Data Collectors Configuration...")
			if err = datacollectors.Import(path.Join(folder, dataCollFileName)); err != nil {
				return err
			}
		}

		clilog.EnablePrintOutput(false)

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
			clilog.Info.Println("Importing configuration for environment " + environment)
			apiclient.SetApigeeEnv(environment)

			if utils.FileExists(path.Join(folder, environment+"_"+keyStoresFileName)) {
				clilog.Info.Println("\tImporting Keystore names...")
				if err = keystores.Import(conn, path.Join(folder, environment+"_"+keyStoresFileName)); err != nil {
					return err
				}
			}

			if utils.FileExists(path.Join(folder, environment+"_"+targetServerFileName)) {
				clilog.Info.Println("\tImporting Target servers...")
				if err = targetservers.Import(conn, path.Join(folder, environment+"_"+targetServerFileName)); err != nil {
					return err
				}
			}

			if utils.FileExists(path.Join(folder, environment+"_"+referencesFileName)) {
				clilog.Info.Println("\tImporting References...")
				if err = references.Import(conn, path.Join(folder, environment+"_"+referencesFileName)); err != nil {
					return err
				}
			}

			if utils.FileExists(path.Join(folder, "env_"+environment+"_"+kVMFileName)) {
				clilog.Info.Println("\tImporting KVM Names only...")
				if kvmList, err = utils.ReadEntityFile(path.Join(folder, "env_"+environment+"_"+kVMFileName)); err != nil {
					return err
				}
				for _, kvmName := range kvmList {
					//create only encrypted KVMs
					if _, err = kvm.Create("", kvmName, true); err != nil {
						return err
					}
					if envKVMFileList[kvmName] != "" {
						if err = kvm.ImportEntries("", kvmName, conn, envKVMFileList[kvmName]); err != nil {
							return err
						}
					}
				}
			}

			if importDebugmask {
				if utils.FileExists(path.Join(folder, environment+debugmaskFileName)) {
					clilog.Info.Println("\tImporting Debug Mask configuration...")
					debugMask, _ := readEntityFileAsString(path.Join(folder, environment+debugmaskFileName))
					if _, err = env.SetDebug(debugMask); err != nil {
						return err
					}
				}
			}

			if importTrace {
				if utils.FileExists(path.Join(folder, environment+tracecfgFileName)) {
					clilog.Info.Println("\tImporting Trace configuration...")
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

func readEntityFileAsString(filePath string) (string, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	return string(byteValue[:]), nil
}

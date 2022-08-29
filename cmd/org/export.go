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
	"os"
	"strconv"
	"strings"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/apis"
	"github.com/apigee/apigeecli/client/apps"
	"github.com/apigee/apigeecli/client/datacollectors"
	"github.com/apigee/apigeecli/client/developers"
	"github.com/apigee/apigeecli/client/env"
	"github.com/apigee/apigeecli/client/envgroups"
	"github.com/apigee/apigeecli/client/keystores"
	"github.com/apigee/apigeecli/client/kvm"
	"github.com/apigee/apigeecli/client/orgs"
	"github.com/apigee/apigeecli/client/products"
	"github.com/apigee/apigeecli/client/references"
	"github.com/apigee/apigeecli/client/sharedflows"
	"github.com/apigee/apigeecli/client/sync"
	"github.com/apigee/apigeecli/client/targetservers"
	"github.com/apigee/apigeecli/clilog"
	"github.com/spf13/cobra"
)

// ExportCmd to get org details
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Apigee Configuration",
	Long:  "Export Apigee Configuration",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		var productResponse, appsResponse, targetServerResponse, referencesResponse [][]byte
		var respBody, listKVMBytes []byte

		runtimeType, _ := orgs.GetOrgField("runtimeType")

		if err = createFolders(); proceedOnError(err) != nil {
			return err
		}

		clilog.Warning.Println("Calls to Apigee APIs have a quota of 6000 per min. Running this tool against large list of entities can exhaust that quota and impact the usage of the platform.")

		fmt.Println("Exporting API Proxies...")
		if err = apis.ExportProxies(conn, proxiesFolderName, allRevisions); proceedOnError(err) != nil {
			return err
		}

		fmt.Println("Exporting Sharedflows...")
		if err = sharedflows.Export(conn, sharedFlowsFolderName, allRevisions); proceedOnError(err) != nil {
			return err
		}

		fmt.Println("Exporting API Products...")
		if productResponse, err = products.Export(conn); proceedOnError(err) != nil {
			return err
		}
		if err = apiclient.WriteArrayByteArrayToFile(productsFileName, false, productResponse); proceedOnError(err) != nil {
			return err
		}

		fmt.Printf("\tExporting KV Map names for org %s\n", org)
		if listKVMBytes, err = kvm.List(""); proceedOnError(err) != nil {
			return err
		}
		if err = apiclient.WriteByteArrayToFile(org+"_"+kVMFileName, false, respBody); proceedOnError(err) != nil {
			return err
		}

		if err = exportKVMEntries("org", "", listKVMBytes); proceedOnError(err) != nil {
			return err
		}

		fmt.Println("Exporting Developers...")
		if respBody, err = developers.Export(); proceedOnError(err) != nil {
			return err
		}
		if err = apiclient.WriteByteArrayToFile(developersFileName, false, respBody); proceedOnError(err) != nil {
			return err
		}

		fmt.Println("Exporting Developer Apps...")
		if appsResponse, err = apps.Export(conn); proceedOnError(err) != nil {
			return err
		}
		if err = apiclient.WriteArrayByteArrayToFile(appsFileName, false, appsResponse); proceedOnError(err) != nil {
			return err
		}

		fmt.Println("Exporting Environment Group Configuration...")
		apiclient.SetPrintOutput(false)
		if respBody, err = envgroups.List(); proceedOnError(err) != nil {
			return err
		}
		if err = apiclient.WriteByteArrayToFile(envGroupsFileName, false, respBody); proceedOnError(err) != nil {
			return err
		}

		fmt.Println("Exporting Data collectors Configuration...")
		if respBody, err = datacollectors.List(); proceedOnError(err) != nil {
			return err
		}
		if err = apiclient.WriteByteArrayToFile(dataCollFileName, false, respBody); proceedOnError(err) != nil {
			return err
		}

		if runtimeType == "HYBRID" {
			fmt.Println("Exporting Sync Authorization Identities...")
			if respBody, err = sync.Get(); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(syncAuthFileName, false, respBody); proceedOnError(err) != nil {
				return err
			}
		}

		var envRespBody []byte
		if envRespBody, err = env.List(); proceedOnError(err) != nil {
			return err
		}

		environments := []string{}
		if err = json.Unmarshal(envRespBody, &environments); proceedOnError(err) != nil {
			return err

		}

		for _, environment := range environments {
			fmt.Println("Exporting configuration for environment " + environment)
			apiclient.SetApigeeEnv(environment)
			fmt.Println("\tExporting Target servers...")
			if targetServerResponse, err = targetservers.Export(conn); proceedOnError(err) != nil {
				return err
			}
			if err = apiclient.WriteArrayByteArrayToFile(environment+"_"+targetServerFileName, false, targetServerResponse); proceedOnError(err) != nil {
				return err
			}

			fmt.Printf("\tExporting KV Map names for environment %s...\n", environment)
			if listKVMBytes, err = kvm.List(""); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(environment+"_"+kVMFileName, false, respBody); proceedOnError(err) != nil {
				return err
			}

			if err = exportKVMEntries("env", environment, listKVMBytes); proceedOnError(err) != nil {
				return err
			}

			fmt.Println("\tExporting Key store names...")
			if respBody, err = keystores.List(); proceedOnError(err) != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(environment+"_"+keyStoresFileName, false, respBody); proceedOnError(err) != nil {
				return err
			}

			fmt.Println("\tExporting debugmask configuration...")
			if respBody, err = env.GetDebug(); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(environment+"_"+debugmaskFileName, false, respBody); proceedOnError(err) != nil {
				return err
			}

			fmt.Println("\tExporting traceconfig...")
			if respBody, err = env.GetTraceConfig(); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(environment+"_"+tracecfgFileName, false, respBody); proceedOnError(err) != nil {
				return err
			}

			fmt.Println("\tExporting references...")
			if referencesResponse, err = references.Export(conn); proceedOnError(err) != nil {
				return err
			}
			if err = apiclient.WriteArrayByteArrayToFile(environment+"_"+referencesFileName, false, referencesResponse); proceedOnError(err) != nil {
				return err
			}

		}

		return
	},
}

var allRevisions, continueOnErr bool

func init() {

	ExportCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	ExportCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	ExportCmd.Flags().BoolVarP(&allRevisions, "all", "",
		false, "Export all revisions, default=false. Exports the latest revision")
	ExportCmd.Flags().BoolVarP(&continueOnErr, "continueOnError", "",
		false, "Ignore errors and continue exporting data")
}

func createFolders() (err error) {
	if err = os.Mkdir(proxiesFolderName, 0755); err != nil {
		return err
	}
	if err = os.Mkdir(sharedFlowsFolderName, 0755); err != nil {
		return err
	}
	return nil
}

func exportKVMEntries(scope string, env string, listKVMBytes []byte) (err error) {

	var kvmEntries [][]byte
	var listKVM []string
	var fileName string

	if err = json.Unmarshal(listKVMBytes, &listKVM); err != nil {
		return err
	}

	for _, mapName := range listKVM {

		fmt.Printf("\tExporting KVM entries for %s in org %s\n", org, mapName)
		if kvmEntries, err = kvm.ExportEntries("", mapName); err != nil {
			return err
		}

		if scope == "org" {
			fileName = strings.Join([]string{scope, mapName, "kvmfile"}, "_")
		} else if scope == "env" {
			fileName = strings.Join([]string{scope, env, mapName, "kvmfile"}, "_")
		}

		if len(kvmEntries) > 0 {
			for i := range kvmEntries {
				if err = apiclient.WriteByteArrayToFile(fileName+"_"+strconv.Itoa(i)+".json", false, kvmEntries[i]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func proceedOnError(e error) error {
	if continueOnErr {
		return nil
	}
	return e
}

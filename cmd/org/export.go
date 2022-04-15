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

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/apis"
	"github.com/srinandan/apigeecli/client/apps"
	"github.com/srinandan/apigeecli/client/datacollectors"
	"github.com/srinandan/apigeecli/client/developers"
	"github.com/srinandan/apigeecli/client/env"
	"github.com/srinandan/apigeecli/client/envgroups"
	"github.com/srinandan/apigeecli/client/keystores"
	"github.com/srinandan/apigeecli/client/kvm"
	"github.com/srinandan/apigeecli/client/orgs"
	"github.com/srinandan/apigeecli/client/products"
	"github.com/srinandan/apigeecli/client/references"
	"github.com/srinandan/apigeecli/client/sharedflows"
	"github.com/srinandan/apigeecli/client/sync"
	"github.com/srinandan/apigeecli/client/targetservers"
	"github.com/srinandan/apigeecli/clilog"
)

//ExportCmd to get org details
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Apigee Configuration",
	Long:  "Export Apigee Configuration",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		var productResponse, appsResponse, targetServerResponse [][]byte
		var respBody []byte

		runtimeType, _ := orgs.GetOrgField("runtimeType")

		if err = createFolders(); err != nil {
			return err
		}

		clilog.Warning.Println("Calls to Apigee APIs have a quota of 6000 per min. Running this tool against large list of entities can exhaust that quota and impact the usage of the platform.")

		fmt.Println("Exporting API Proxies...")
		if err = apis.ExportProxies(conn, proxiesFolderName, allRevisions); err != nil {
			return err
		}

		fmt.Println("Exporting Sharedflows...")
		if err = sharedflows.Export(conn, sharedFlowsFolderName, allRevisions); err != nil {
			return err
		}

		fmt.Println("Exporting API Products...")
		if productResponse, err = products.Export(conn); err != nil {
			return err
		}
		if err = apiclient.WriteArrayByteArrayToFile(productsFileName, false, productResponse); err != nil {
			return err
		}

		fmt.Println("\tExporting KV Map names for org %s", org)
		if respBody, err = kvm.List(""); err != nil {
			return err
		}
		if err = apiclient.WriteByteArrayToFile(org+"_"+kVMFileName, false, respBody); err != nil {
			return err
		}

		fmt.Println("Exporting Developers...")
		if respBody, err = developers.Export(); err != nil {
			return err
		}
		if err = apiclient.WriteByteArrayToFile(developersFileName, false, respBody); err != nil {
			return err
		}

		fmt.Println("Exporting Developer Apps...")
		if appsResponse, err = apps.Export(conn); err != nil {
			return err
		}
		if err = apiclient.WriteArrayByteArrayToFile(appsFileName, false, appsResponse); err != nil {
			return err
		}

		fmt.Println("Exporting Environment Group Configuration...")
		apiclient.SetPrintOutput(false)
		if respBody, err = envgroups.List(); err != nil {
			return err
		}
		if err = apiclient.WriteByteArrayToFile(envGroupsFileName, false, respBody); err != nil {
			return err
		}

		fmt.Println("Exporting Data collectors Configuration...")
		if respBody, err = datacollectors.List(); err != nil {
			return err
		}
		if err = apiclient.WriteByteArrayToFile(dataCollFileName, false, respBody); err != nil {
			return err
		}

		if runtimeType == "HYBRID" {
			fmt.Println("Exporting Sync Authorization Identities...")
			if respBody, err = sync.Get(); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(syncAuthFileName, false, respBody); err != nil {
				return err
			}
		}

		var envRespBody []byte
		if envRespBody, err = env.List(); err != nil {
			return err
		}

		environments := []string{}
		if err = json.Unmarshal(envRespBody, &environments); err != nil {
			return err

		}

		for _, environment := range environments {
			fmt.Println("Exporting configuration for environment " + environment)
			apiclient.SetApigeeEnv(environment)
			fmt.Println("\tExporting Target servers...")
			if targetServerResponse, err = targetservers.Export(conn); err != nil {
				return err
			}
			if err = apiclient.WriteArrayByteArrayToFile(environment+"_"+targetServerFileName, false, targetServerResponse); err != nil {
				return err
			}

			fmt.Println("\tExporting References...")
			if respBody, err = references.List(); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(environment+"_"+refFileName, false, respBody); err != nil {
				return err
			}

			fmt.Println("\tExporting KV Map names for environment %s...", environment)
			if respBody, err = kvm.List(""); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(environment+"_"+kVMFileName, false, respBody); err != nil {
				return err
			}

			fmt.Println("\tExporting Key store names...")
			if respBody, err = keystores.List(); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(environment+"_"+kVMFileName, false, respBody); err != nil {
				return err
			}

			fmt.Println("\tExporting debugmask configuration...")
			if respBody, err = env.GetDebug(); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(environment+debugmaskFileName, false, respBody); err != nil {
				return err
			}

			fmt.Println("\tExporting traceconfig...")
			if respBody, err = env.GetTraceConfig(); err != nil {
				return err
			}
			if err = apiclient.WriteByteArrayToFile(environment+tracecfgFileName, false, respBody); err != nil {
				return err
			}

		}

		return
	},
}

var allRevisions bool

func init() {

	ExportCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	ExportCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
	ExportCmd.Flags().BoolVarP(&allRevisions, "all", "",
		false, "Export all revisions, default=false. Exports the latest revision")
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

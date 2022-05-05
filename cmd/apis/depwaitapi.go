// Copyright 2020 Google LLC
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

package apis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/apis"
	"github.com/spf13/cobra"
)

//DepWaitCmd to deploy api
var DepWaitCmd = &cobra.Command{
	Use:   "deploy-wait",
	Short: "Deploys a revision of an existing API proxy and waits for deployment status",
	Long:  "Deploys a revision of an existing API proxy to an environment and waits for deployment status",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		apiclient.SetApigeeEnv(env)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if _, err = apis.DeployProxy(name, revision, overrides, serviceAccountName); err != nil {
			return err
		}

		fmt.Printf("Checking deployment status in %d seconds\n", interval)

		apiclient.SetPrintOutput(false)

		stop := apiclient.Every(interval*time.Second, func(time.Time) bool {
			var respBody []byte
			respMap := make(map[string]interface{})
			if respBody, err = apis.ListProxyRevisionDeployments(name, revision); err != nil {
				return true
			}

			if err = json.Unmarshal(respBody, &respMap); err != nil {
				return true
			}

			if respMap["state"] == "PROGRESSING" {
				fmt.Printf("Proxy deployment status is: %s. Waiting %d seconds.\n", respMap["state"], interval)
				return true
			} else if respMap["state"] == "READY" {
				fmt.Println("Proxy deployment completed with status: ", respMap["state"])
				return false
			} else {
				fmt.Println("Proxy deployment failed with status: ", respMap["state"])
				return false
			}
		})

		<-stop

		return
	},
}

const interval = 10

func init() {

	DepWaitCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	DepWaitCmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	DepWaitCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision")
	DepWaitCmd.Flags().BoolVarP(&overrides, "ovr", "r",
		false, "Forces deployment of the new revision")
	DepWaitCmd.Flags().StringVarP(&serviceAccountName, "sa", "s",
		"", "The format must be {ACCOUNT_ID}@{PROJECT}.iam.gserviceaccount.com.")

	_ = DepWaitCmd.MarkFlagRequired("env")
	_ = DepWaitCmd.MarkFlagRequired("name")
}

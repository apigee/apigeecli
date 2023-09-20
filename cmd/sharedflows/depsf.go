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

package sharedflows

import (
	"encoding/json"
	"time"

	"internal/apiclient"
	"internal/clilog"

	"internal/client/sharedflows"

	"github.com/spf13/cobra"
)

// DepCmd to deploy shared flow
var DepCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a revision of an existing Sharedflow",
	Long:  "Deploys a revision of an existing Sharedflow to an environment in an organization",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if revision == -1 {
			if revision, err = sharedflows.GetHighestSfRevision(name); err != nil {
				return err
			}
		}
		_, err = sharedflows.Deploy(name, revision, overrides, serviceAccountName)
		apiclient.DisableCmdPrintHttpResponse()

		if wait {
			clilog.Info.Printf("Checking deployment status in %d seconds\n", interval)

			stop := apiclient.Every(interval*time.Second, func(time.Time) bool {
				var respBody []byte
				respMap := make(map[string]interface{})
				if respBody, err = sharedflows.ListRevisionDeployments(name, revision); err != nil {
					clilog.Error.Printf("Error fetching sharedflow revision status: %v", err)
					return false
				}

				if err = json.Unmarshal(respBody, &respMap); err != nil {
					return true
				}

				if respMap["state"] == "PROGRESSING" {
					clilog.Info.Printf("Sharedflow deployment status is: %s. Waiting %d seconds.\n", respMap["state"], interval)
					return true
				} else if respMap["state"] == "READY" {
					clilog.Info.Println("Sharedflow deployment completed with status: ", respMap["state"])
					return false
				} else {
					clilog.Info.Println("Sharedflow deployment failed with status: ", respMap["state"])
					return false
				}
			})

			<-stop
		}

		return err
	},
}

var (
	overrides, wait    bool
	serviceAccountName string
)

const interval = 10

func init() {
	DepCmd.Flags().StringVarP(&name, "name", "n",
		"", "Sharedflow name")
	DepCmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	DepCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "Sharedflow revision. If not set, the highest revision is used")
	DepCmd.Flags().BoolVarP(&overrides, "ovr", "r",
		false, "Forces deployment of the new revision")
	DepCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the deployment to finish, with success or error")

	DepCmd.Flags().StringVarP(&serviceAccountName, "sa", "s",
		"", "The format must be {ACCOUNT_ID}@{PROJECT}.iam.gserviceaccount.com.")

	_ = DepCmd.MarkFlagRequired("env")
	_ = DepCmd.MarkFlagRequired("name")
}

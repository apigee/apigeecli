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

package overrides

import (
	"internal/apiclient"

	"internal/clilog"

	"internal/client/env"
	"internal/client/envgroups"

	"github.com/spf13/cobra"
)

// DeleteCmd provisions control plane entities for hybrid
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete control plane entities",
	Long:  "Delete control plane entities",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = readOverrides(overridesFile); err != nil {
			return err
		}
		apiclient.SetProjectID(getOrg())
		apiclient.DisableCmdPrintHttpResponse()
		_ = apiclient.SetApigeeOrg(getOrg())

		// delete environments
		environmentList := getEnvs()
		for _, environment := range environmentList {
			// check if env exists
			apiclient.SetApigeeEnv(environment)
			if _, err = env.Get(false); err != nil {
				if _, err = env.Delete(); err != nil {
					return err
				}
				clilog.Info.Printf("Environment %s deleted", environment)
			} else {
				clilog.Info.Printf("Environment %s does not exist\n", environment)
			}
		}

		// create environment groups
		environmentGroupList := getEnvGroups()
		for _, environmentGroup := range environmentGroupList {
			// check if env group exists
			if _, err = envgroups.Get(environmentGroup); err != nil {
				if _, err = envgroups.Delete(environmentGroup); err != nil {
					return err
				}
				clilog.Info.Printf("Environment Group %s deleted\n", environmentGroup)
			} else {
				clilog.Info.Printf("Environment Group %s does not exist\n", environmentGroup)
			}
		}

		return
	},
}

func init() {
	DeleteCmd.Flags().StringVarP(&overridesFile, "overrides", "f",
		"overrides.yaml", "overrides file path")

	_ = DeleteCmd.MarkFlagRequired("overrides")
}

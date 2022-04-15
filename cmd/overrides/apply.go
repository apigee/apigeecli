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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/env"
	"github.com/srinandan/apigeecli/client/envgroups"
	"github.com/srinandan/apigeecli/client/orgs"
	"github.com/srinandan/apigeecli/client/sync"
	"github.com/srinandan/apigeecli/clilog"
)

//ApplyCmd provisions control plane entities for hybrid
var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply control plane entities",
	Long:  "Apply control plane entities",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if err = readOverrides(overridesFile); err != nil {
			return err
		}
		apiclient.SetProjectID(getOrg())
		apiclient.SetApigeeOrg(getOrg())
		apiclient.SetPrintOutput(false)

		//check if the org exists
		if _, err = orgs.Get(); err != nil {
			if _, err = orgs.Create(getOrgRegion(), "", "HYBRID", "", "", false); err != nil {
				return err
			}
			fmt.Printf("Org %s created\n", getOrg())
		} else {
			fmt.Printf("Org %s already exists\n", getOrg())
		}

		//check setSyncAuth
		identities := getSyncServiceAccounts()
		if len(identities) > 0 {
			if _, err = sync.Set(identities); err != nil {
				clilog.Warning.Println("Error setting identities: ", err)
			} else {
				fmt.Printf("Org setSync identities set: %v", identities)
			}
		} else {
			clilog.Warning.Println("No sync identities were set in overrides")
		}

		//create environments
		environmentList := getEnvs()
		for _, environment := range environmentList {
			//check if env exists
			apiclient.SetApigeeEnv(environment)
			if _, err = env.Get(false); err != nil {
				if _, err = env.Create("PROXY", "PROGRAMMABLE"); err != nil {
					return err
				}
				fmt.Printf("Environment %s created", environment)
			} else {
				fmt.Printf("Environment %s already exists\n", environment)
			}
		}

		//create environment groups
		environmentGroupList := getEnvGroups()
		for i, environmentGroup := range environmentGroupList {
			//check if env group exists
			if _, err = envgroups.Get(environmentGroup); err != nil {
				if _, err = envgroups.Create(environmentGroup, getDomainName(i)); err != nil {
					return err
				}
				fmt.Printf("Environment Group %s provisioned with a temporary domain name %s\n", environmentGroup, getDomainName(i))
			} else {
				fmt.Printf("Environment Group %s already exists\n", environmentGroup)
			}
		}

		return
	},
}

var overridesFile string

func init() {

	ApplyCmd.Flags().StringVarP(&overridesFile, "overrides", "f",
		"overrides.yaml", "overrides file path")

	_ = ApplyCmd.MarkFlagRequired("overrides")
}

func getDomainName(index int) []string {
	var domainNames = []string{}
	domainName := fmt.Sprintf("api.acme%d.com", index)
	domainNames = append(domainNames, domainName)
	return domainNames
}

func getSyncIdentityCount(response []byte) (count int, err error) {
	syncResponse := make(map[string]interface{})
	err = json.Unmarshal(response, &syncResponse)
	if err != nil {
		return 0, err
	}
	return len(syncResponse["identities"].([]interface{})), nil
}

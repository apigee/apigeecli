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

package org

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/orgs"
)

//Cmd to set mart endpoint
var SetAddonCmd = &cobra.Command{
	Use:   "setaddons",
	Short: "Enable addons for an Apigee organization",
	Long:  "Enable addons for an Apigee organization",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = orgs.SetAddons(advancedApiOpsConfig, integrationConfig, monetizationConfig)
		return
	},
}

var advancedApiOpsConfig, integrationConfig, monetizationConfig bool

func init() {

	SetAddonCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	SetAddonCmd.Flags().BoolVarP(&advancedApiOpsConfig, "enable-advapiops", "",
		false, "Enable Apigee Advanced API OPs")
	SetAddonCmd.Flags().BoolVarP(&integrationConfig, "enable-integration", "",
		false, "Enable Apigee Integration")
	SetAddonCmd.Flags().BoolVarP(&monetizationConfig, "enable-mint", "",
		false, "Enable Apigee Monetization")

}

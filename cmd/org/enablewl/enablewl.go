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

package enablewl

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/orgs"
)

//Cmd to set mart endpoint
var Cmd = &cobra.Command{
	Use:   "enable-mart-whitelist",
	Short: "Enable IP whitelisting for MART connections",
	Long:  "Enable IP whitelisting for MART connections",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return orgs.SetOrgProperty("features.mart.ip.whitelist.enabled", "true")
	},
}

func init() {

	Cmd.Flags().StringVarP(apiclient.GetApigeeOrgP(), "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkFlagRequired("org")
}

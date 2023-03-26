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
	"internal/apiclient"

	"internal/client/apis"

	"github.com/spf13/cobra"
)

// CreateTrcCmd to manage tracing of apis
var CreateTrcCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new debug session for an API proxy",
	Long:  "Create a new debug session for Apigee API proxy revision deployed in an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.CreateTraceSession(name, revision, filter)
		return
	},
}

var filter map[string]string

func init() {

	CreateTrcCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	CreateTrcCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision")
	CreateTrcCmd.Flags().StringToStringVar(&filter, "filter",
		nil, "Filter Conditions; format is name1=value1,name2=value2...")

	_ = CreateTrcCmd.MarkFlagRequired("name")
	_ = CreateTrcCmd.MarkFlagRequired("rev")
}

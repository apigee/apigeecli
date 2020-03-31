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

package env

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/env"
)

//Cmd to enable AX obfuscation
var ObCmd = &cobra.Command{
	Use:   "enable-ax-obfuscation",
	Short: "Obfuscate analytics fields",
	Long:  "Obfuscate analytics fields before sending to control plane",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		apiclient.SetApigeeEnv(environment)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return env.SetEnvProperty("features.analytics.data.obfuscation.enabled", "true")
	},
}

func init() {

	ObCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")

	_ = ObCmd.MarkFlagRequired("env")
}

// Copyright 2024 Google LLC
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

package kvm

import (
	"fmt"
	"internal/apiclient"
	"internal/client/kvm"

	"github.com/spf13/cobra"
)

// UpdateCmd to update a kvm
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a KV Map",
	Long:  "Update a KV Map",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if env != "" {
			apiclient.SetApigeeEnv(env)
		}
		if env != "" && proxyName != "" {
			return fmt.Errorf("proxy and env flags cannot be used together")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = kvm.Update(proxyName, name, masked)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&env, "env", "e",
		"", "Environment name")
	UpdateCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API Proxy name")
	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "KVM Map name")
	UpdateCmd.Flags().BoolVarP(&masked, "masked", "m",
		false, "Mask the KVM values in the debug session")

	_ = UpdateCmd.MarkFlagRequired("name")
}

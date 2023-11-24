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
	"internal/apiclient"

	environments "internal/client/env"

	"github.com/spf13/cobra"
)

// EnableSecActCmd to manage tracing of apis
var EnableSecActCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a SecurityAction",
	Long:  "Enable a SecurityAction",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = environments.EnableSecurityAction(name)
		return
	},
}

func init() {
	EnableSecActCmd.Flags().StringVarP(&name, "name", "n",
		"", "Security Action name")
	_ = EnableSecActCmd.MarkFlagRequired("name")
}
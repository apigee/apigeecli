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

	"internal/client/env"

	"github.com/spf13/cobra"
)

// PropCmd to set env prop
var PropCmd = &cobra.Command{
	Use:   "set",
	Short: "Set environment property",
	Long:  "Set environment property",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return env.SetEnvProperty(propName, propValue)
	},
}

var propName, propValue string

func init() {
	PropCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")
	PropCmd.Flags().StringVarP(&propName, "name", "n",
		"", "Property name")
	PropCmd.Flags().StringVarP(&propValue, "value", "v",
		"", "Property Value")

	_ = PropCmd.MarkFlagRequired("env")
	_ = PropCmd.MarkFlagRequired("name")
	_ = PropCmd.MarkFlagRequired("value")
}

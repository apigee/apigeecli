// Copyright 2023 Google LLC
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
)

// SecActCmd to manage security incidents
var SecActCmd = &cobra.Command{
	Use:   "secactions",
	Short: "Manage SecurityActions for Apigee Advanced Security",
	Long:  "Manage SecurityActions for Apigee Advanced Security",
}

func init() {
	SecActCmd.PersistentFlags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")

	_ = SecActCmd.MarkPersistentFlagRequired("env")

	SecActCmd.AddCommand(ListSecActCmd)
	SecActCmd.AddCommand(GetSecActCmd)
	SecActCmd.AddCommand(EnableSecActCmd)
	SecActCmd.AddCommand(DisableSecActCmd)
}

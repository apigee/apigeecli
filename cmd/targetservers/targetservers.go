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

package targetservers

import (
	"github.com/spf13/cobra"
)

// Cmd to manage targetservers
var Cmd = &cobra.Command{
	Use:     "targetservers",
	Aliases: []string{"ts"},
	Short:   "Manage Target Servers",
	Long:    "Manage Target Servers",
}

var org, env, name string

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(ExpCmd)
	Cmd.AddCommand(ImpCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(UpdateCmd)
}

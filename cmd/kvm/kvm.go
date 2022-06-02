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

package kvm

import (
	"github.com/spf13/cobra"
)

//Cmd to manage kvms
var Cmd = &cobra.Command{
	Use:   "kvms",
	Short: "Manage Key Value Maps",
	Long:  "Manage Key Value Maps",
}

var org, env, name, proxyName string

func init() {

	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(EntryCmd)
}

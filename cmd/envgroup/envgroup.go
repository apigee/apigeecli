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

package envgroup

import (
	"github.com/spf13/cobra"
)

// Cmd to manage envs
var Cmd = &cobra.Command{
	Use:   "envgroups",
	Short: "Manage Apigee environment groups",
	Long:  "Manage Apigee environment groups",
}

var org, name, environment string
var hostnames []string

func init() {

	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	Cmd.AddCommand(AttachCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(ListAttachCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(DetachCmd)
	Cmd.AddCommand(UpdateCmd)
	Cmd.AddCommand(DelCmd)
}

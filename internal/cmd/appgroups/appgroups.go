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

package appgroups

import (
	"github.com/spf13/cobra"
)

// Cmd to manage appgroups
var Cmd = &cobra.Command{
	Use:   "appgroups",
	Short: "Manage Apigee Application Groups",
	Long:  "Manage Apigee Application Groups",
}

var org, region string

var examples = []string{"apigeecli appgroups import -f samples/appgroups.json --default-token"}

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "Apigee control plane region name; default is https://apigee.googleapis.com")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(UpdateCmd)
	Cmd.AddCommand(ManageCmd)
	Cmd.AddCommand(AppCmd)
	Cmd.AddCommand(ExpCmd)
	Cmd.AddCommand(ImpCmd)
}

func GetExample(i int) string {
	return examples[i]
}

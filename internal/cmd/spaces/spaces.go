// Copyright 2025 Google LLC
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

package spaces

import (
	"github.com/spf13/cobra"
)

// Cmd to manage spaces
var Cmd = &cobra.Command{
	Use:     "spaces",
	Aliases: []string{"space"},
	Short:   "Manage Apigee Spaces",
	Long:    "Manage Apigee Spaces",
}

var org, displayName, name, region string

var examples = []string{"apigeecli spaces create --name=space1 --display-name=\"Space 1\""}

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "Apigee control plane region name; default is https://apigee.googleapis.com")

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(UpdateCmd)
	// TODO Cmd.AddCommand(ExpCmd)
	// TODO Cmd.AddCommand(ImpCmd)
}

func GetExample(i int) string {
	return examples[i]
}

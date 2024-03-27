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

package datacollectors

import (
	"github.com/spf13/cobra"
)

// Cmd to manage datacollectors
var Cmd = &cobra.Command{
	Use:     "datacollectors",
	Aliases: []string{"dc"},
	Short:   "Manage Apigee datacollectors entities",
	Long:    "Manage Apigee datacollectors entities",
}

var org, name, region string

var examples = []string{"apigeecli datacollectors import -f samples/datacollectors.json --default-token"}

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "Apigee control plane region name; default is https://apigee.googleapis.com")

	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(ImpCmd)
	Cmd.AddCommand(ExpCmd)
}

func GetExample(i int) string {
	return examples[i]
}

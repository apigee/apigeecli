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

package sharedflows

import (
	"github.com/spf13/cobra"
)

// Cmd to manage shared flows
var Cmd = &cobra.Command{
	Use:   "sharedflows",
	Short: "Manage Apigee shared flows in an org",
	Long:  "Manage Apigee shared flows in an org",
}

var (
	name, org, env, region, space string
	conn, revision                int
)

var examples = []string{"apigeecli sharedflows import -f samples/sharedflows"}

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "Apigee control plane region name; default is https://apigee.googleapis.com")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(FetCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(UndepCmd)
	Cmd.AddCommand(DepCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(ExpCmd)
	Cmd.AddCommand(ImpCmd)
	Cmd.AddCommand(CleanCmd)
	Cmd.AddCommand(ListDepCmd)
	Cmd.AddCommand(MoveCmd)
}

func GetExample(i int) string {
	return examples[i]
}

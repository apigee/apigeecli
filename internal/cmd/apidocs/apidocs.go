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

package apidocs

import (
	"github.com/spf13/cobra"
)

// Cmd to manage apis
var Cmd = &cobra.Command{
	Use:   "apidocs",
	Short: "Manage Apigee API catalog item through ApiDoc",
	Long:  "Manage Apigee API catalog item through ApiDoc",
}

var org, siteid, id, name, region string

var examples = []string{"apigeecli apidocs import -f samples/apidocs  -s $siteId"}

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&siteid, "siteid", "s",
		"", "Name or siteid of the portal")
	Cmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "Apigee control plane region name; default is https://apigee.googleapis.com")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(DocCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(UpdateCmd)
	Cmd.AddCommand(ExpCmd)
	Cmd.AddCommand(ImpCmd)

	_ = Cmd.MarkFlagRequired("org")
	_ = Cmd.MarkFlagRequired("siteid")
}

func GetExample(i int) string {
	return examples[i]
}

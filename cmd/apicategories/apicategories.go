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

package apicategories

import (
	"github.com/spf13/cobra"
)

// Cmd to manage api catalog items
var Cmd = &cobra.Command{
	Use:   "apicategories",
	Short: "Manage Apigee API categories that are tagged on catalog items",
	Long:  "Manage Apigee API categories that are tagged on catalog items",
}

var org, siteid, id string

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&siteid, "siteid", "s",
		"", "Name or siteid of the portal")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(CreateCmd)

	_ = Cmd.MarkFlagRequired("org")
	_ = Cmd.MarkFlagRequired("siteid")
}

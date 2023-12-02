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

package sites

import (
	"github.com/spf13/cobra"
)

// Cmd to manage
var Cmd = &cobra.Command{
	Use:   "sites",
	Short: "Manage Apigee API Portals",
	Long:  "Manage Apigee API Portals",
}

var org, siteid, id, name string

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	Cmd.AddCommand(ListCmd)

	_ = Cmd.MarkFlagRequired("org")
}

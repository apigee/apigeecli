// Copyright 2024 Google LLC
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

package specs

import (
	"github.com/spf13/cobra"
)

// SpecsCmd to manage apis
var SpecsCmd = &cobra.Command{
	Use:   "specs",
	Short: "Manage Specs for an API Version in API Hub",
	Long:  "Manage Specs for an API Version in API Hub",
}

var org, region string

func init() {
	SpecsCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	SpecsCmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "API Hub region name")

	SpecsCmd.AddCommand(CrtCmd)
	SpecsCmd.AddCommand(GetCmd)
	SpecsCmd.AddCommand(DelCmd)
	SpecsCmd.AddCommand(ListCmd)

	_ = SpecsCmd.MarkFlagRequired("org")
	_ = SpecsCmd.MarkFlagRequired("region")
}

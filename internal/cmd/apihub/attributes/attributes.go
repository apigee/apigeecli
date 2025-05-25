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

package attributes

import (
	"github.com/spf13/cobra"
)

// AttributeCmd to manage apis
var AttributeCmd = &cobra.Command{
	Use:   "attributes",
	Short: "Manage attributes in API Hub",
	Long:  "Manage attributes in API Hub",
}

var org, region string

func init() {
	AttributeCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	AttributeCmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "API Hub region name")

	AttributeCmd.AddCommand(CrtCmd)
	AttributeCmd.AddCommand(GetCmd)
	AttributeCmd.AddCommand(DelCmd)
	AttributeCmd.AddCommand(ListCmd)
	AttributeCmd.AddCommand(UpdateCmd)

	_ = AttributeCmd.MarkFlagRequired("org")
	_ = AttributeCmd.MarkFlagRequired("region")
}

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

package instances

import (
	"github.com/spf13/cobra"
)

// InstanceCmd to manage apis
var InstanceCmd = &cobra.Command{
	Use:   "instances",
	Short: "Manage Apigee API Hub Instances",
	Long:  "Manage Apigee API Hub Instances",
}

var org, region string

func init() {
	InstanceCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	InstanceCmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "API Hub region name")

	InstanceCmd.AddCommand(CrtCmd)
	InstanceCmd.AddCommand(GetCmd)
	InstanceCmd.AddCommand(ListCmd)

	_ = InstanceCmd.MarkFlagRequired("org")
	_ = InstanceCmd.MarkFlagRequired("region")
}

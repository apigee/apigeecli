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

package externalapis

import (
	"github.com/spf13/cobra"
)

// ExternalAPICmd to manage apis
var ExternalAPICmd = &cobra.Command{
	Use:   "externalapis",
	Short: "Manage APIs provided by external sources",
	Long:  "Manage APIs provided by external sources",
}

var org, region string

func init() {
	ExternalAPICmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	ExternalAPICmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "API Hub region name")

	ExternalAPICmd.AddCommand(CrtCmd)
	ExternalAPICmd.AddCommand(GetCmd)
	ExternalAPICmd.AddCommand(DelCmd)
	ExternalAPICmd.AddCommand(ListCmd)
	ExternalAPICmd.AddCommand(UpdateCmd)

	_ = ExternalAPICmd.MarkFlagRequired("org")
	_ = ExternalAPICmd.MarkFlagRequired("region")
}

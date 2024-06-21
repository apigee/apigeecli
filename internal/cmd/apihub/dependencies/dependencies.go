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

package dependencies

import (
	"github.com/spf13/cobra"
)

// DependencyCmd to manage apis
var DependencyCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Manage dependencies between consumers and suppliers",
	Long:  "Manage dependencies between consumers and suppliers",
}

var org, region string

func init() {
	DependencyCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	DependencyCmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "API Hub region name")

	DependencyCmd.AddCommand(CrtCmd)
	DependencyCmd.AddCommand(GetCmd)
	DependencyCmd.AddCommand(DelCmd)
	DependencyCmd.AddCommand(ListCmd)

	_ = DependencyCmd.MarkFlagRequired("org")
	_ = DependencyCmd.MarkFlagRequired("region")
}

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

package specs

import "github.com/spf13/cobra"

// SpecCmd to manage apis
var SpecCmd = &cobra.Command{
	Use:   "specs",
	Short: "Manage specs for an API version",
	Long:  "Manage specs for an API version in Apigee Registry",
}

func init() {
	SpecCmd.AddCommand(ListSpecCmd)
	SpecCmd.AddCommand(GetSpecCmd)
}

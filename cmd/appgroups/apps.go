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

package appgroups

import (
	"github.com/spf13/cobra"
)

// AppCmd to manage apps
var AppCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage apps in an Apigee Application Group",
	Long:  "Manage apps in an Apigee Application Group",
}

func init() {
	AppCmd.AddCommand(CreateAppCmd)
	AppCmd.AddCommand(ListAppCmd)
	AppCmd.AddCommand(GetAppCmd)
	AppCmd.AddCommand(DelAppCmd)
	AppCmd.AddCommand(ManageAppCmd)
	AppCmd.AddCommand(KeyCmd)
	AppCmd.AddCommand(ExpAppCmd)
}

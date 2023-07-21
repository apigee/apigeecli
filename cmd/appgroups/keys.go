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

// KeyCmd to manage keys in an app
var KeyCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage keys in an App within an AppGroup",
	Long:  "Manage keys in an App within an AppGroup",
}

func init() {
	KeyCmd.AddCommand(GetKeyCmd)
	KeyCmd.AddCommand(DelKeyCmd)
	KeyCmd.AddCommand(CreateKeyCmd)
	KeyCmd.AddCommand(ManageKeyCmd)
}

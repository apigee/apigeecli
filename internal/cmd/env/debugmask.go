// Copyright 2020 Google LLC
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

package env

import (
	"github.com/spf13/cobra"
)

// DebugCmd to manage debug masks
var DebugCmd = &cobra.Command{
	Use:   "debugmask",
	Short: "Manage debugmasks for the environment",
	Long:  "Manage debugmasks for the environment",
}

func init() {
	DebugCmd.PersistentFlags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")

	_ = DebugCmd.MarkPersistentFlagRequired("env")

	DebugCmd.AddCommand(GetDebugCmd)
	DebugCmd.AddCommand(SetDebugCmd)
}

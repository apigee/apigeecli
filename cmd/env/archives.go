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

// Cmd to manage tracing of apis
var ArchiveCmd = &cobra.Command{
	Use:   "archives",
	Short: "Manage archive deployments for the environment",
	Long:  "Manage archive deployments for the environment",
}

var name string

func init() {

	ArchiveCmd.PersistentFlags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")

	_ = ArchiveCmd.MarkPersistentFlagRequired("env")

	ArchiveCmd.AddCommand(ListArchiveCmd)
	ArchiveCmd.AddCommand(CreateArchiveCmd)
	ArchiveCmd.AddCommand(GetArchiveCmd)
	ArchiveCmd.AddCommand(DelArchiveCmd)
}

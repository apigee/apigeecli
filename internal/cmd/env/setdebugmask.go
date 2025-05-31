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
	"internal/apiclient"
	"internal/cmd/utils"

	environments "internal/client/env"

	"github.com/spf13/cobra"
)

// SetDebugCmd to set debug mas  on env
var SetDebugCmd = &cobra.Command{
	Use:   "set",
	Short: "Set debugmasks for an Environment",
	Long:  "Set debugmasks for an Environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		payload, err := utils.ReadFile(filePath)
		if err != nil {
			return nil
		}
		_, err = environments.SetDebug(payload)
		return
	},
}

func init() {
	SetDebugCmd.Flags().StringVarP(&filePath, "mask-file", "f",
		"", "A path to a file containing debug mask configuration")

	_ = SetDebugCmd.MarkFlagRequired("mask-file")
}

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

package sharedflows

import (
	"fmt"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/sharedflows"
	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// Cmd to import shared flow
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a folder containing sharedflow bundles",
	Long:  "Import a folder containing sharedflow bundles",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if !utils.TestFolder(folder) {
			return fmt.Errorf("supplied path is not a folder")
		}
		return sharedflows.Import(conn, folder)
	},
}

var folder string

func init() {
	ImpCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "folder containing sharedflow bundles")
	ImpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = ImpCmd.MarkFlagRequired("folder")
}

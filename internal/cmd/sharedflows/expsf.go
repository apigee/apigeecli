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
	"internal/apiclient"
	"internal/client/sharedflows"

	"github.com/spf13/cobra"
)

// ExpCmd to export shared flows
var ExpCmd = &cobra.Command{
	Use:   "export",
	Short: "export Sharedflow bundles from an org",
	Long:  "export Sharedflow bundles from an org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if err = apiclient.FolderExists(folder); err != nil {
			return err
		}
		return sharedflows.Export(conn, folder, allRevisions, space)
	},
}

var allRevisions bool

func init() {
	ExpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
	ExpCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "folder to export sharedflow bundles")
	ExpCmd.Flags().BoolVarP(&allRevisions, "all", "",
		false, "Export all proxy revisions")
	ExpCmd.Flags().StringVarP(&space, "space", "",
		"", "Apigee Space associated to")
}

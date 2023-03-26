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

package apps

import (
	"internal/apiclient"

	"github.com/apigee/apigeecli/client/apps"
	"github.com/spf13/cobra"
)

// ImpCmd to import apps
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing Developer Apps",
	Long:  "Import a file containing Developer Apps",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return apps.Import(conn, filePath, developersFilePath)
	},
}

var filePath, developersFilePath string

func init() {

	ImpCmd.Flags().StringVarP(&filePath, "file", "f",
		"", "File containing Developer Apps")
	ImpCmd.Flags().StringVarP(&developersFilePath, "dev-file", "d",
		"", "File containing Developers")
	ImpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = ImpCmd.MarkFlagRequired("file")
	_ = ImpCmd.MarkFlagRequired("dev-file")
}

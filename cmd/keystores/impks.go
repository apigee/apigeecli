// Copyright 2021 Google LLC
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

package keystores

import (
	"internal/apiclient"

	"internal/client/keystores"

	"github.com/spf13/cobra"
)

// ImpCmd to import ts
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing keystores",
	Long:  "Import a file containing keystores",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return keystores.Import(conn, filePath)
	},
}

var (
	filePath string
	conn     int
)

func init() {
	ImpCmd.Flags().StringVarP(&filePath, "file", "f",
		"", "Path to file containing Apigee Keystores")
	ImpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = ImpCmd.MarkFlagRequired("file")
}

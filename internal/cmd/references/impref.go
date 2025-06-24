// Copyright 2022 Google LLC
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

package references

import (
	"internal/apiclient"
	"internal/client/references"

	"github.com/spf13/cobra"
)

// ImpCmd to import references
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Environment References",
	Long:  "Import Environment References",
	Example: `Import Environment References from a file:
` + GetExample(0),
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		return references.Import(conn, filePath)
	},
}

var filePath string

func init() {
	ImpCmd.Flags().StringVarP(&filePath, "file", "f",
		"", "File containing API Products")
	ImpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = ImpCmd.MarkFlagRequired("file")
}

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

package securityprofiles

import (
	"os"

	"internal/apiclient"

	"internal/client/securityprofiles"

	"github.com/spf13/cobra"
)

// ExpCmd to export sec profiles
var ExpCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Security Profiles to a file",
	Long:  "Export Security Profiles to a file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if folder == "" {
			folder, _ = os.Getwd()
		}
		if err = apiclient.FolderExists(folder); err != nil {
			return err
		}
		return securityprofiles.Export(conn, folder, allRevisions)
	},
}

var (
	conn         int
	folder       string
	allRevisions bool
)

func init() {
	ExpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
	ExpCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "folder to export Security Profiles")
	ExpCmd.Flags().BoolVarP(&allRevisions, "all", "",
		false, "Export all proxy revisions")
}

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

package securityprofiles

import (
	"fmt"
	"internal/apiclient"
	"internal/client/securityprofiles"
	"os"

	"github.com/spf13/cobra"
)

// ImpCmd to import sec profiles
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a folder containing Security Profiles",
	Long:  "Import a folder containing Security Profiles",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if stat, err := os.Stat(folder); err == nil && !stat.IsDir() {
			return fmt.Errorf("supplied path is not a folder")
		}
		return securityprofiles.Import(conn, folder)
	},
}

func init() {
	ImpCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "folder containing one or more security profiles")
	ImpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = ImpCmd.MarkFlagRequired("folder")
}

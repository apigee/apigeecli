// Copyright 2024 Google LLC
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

package sites

import (
	"fmt"
	"internal/apiclient"
	"internal/client/sites"
	"os"

	"github.com/spf13/cobra"
)

// CreateCmd creates an integrated developer portal
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Apigee integrated developer portal",
	Long:  "Create an Apigee integrated developer portal from a JSON file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if createFilePath == "" {
			return fmt.Errorf("required flag \"file\" not set")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		content, err := os.ReadFile(createFilePath)
		if err != nil {
			return err
		}
		_, err = sites.Create(createSiteID, content)
		return
	},
}

var (
	createSiteID   string
	createFilePath string
)

func init() {
	CreateCmd.Flags().StringVarP(&createSiteID, "site", "s", "", "Portal site ID (optional, auto-generated if not set)")
	CreateCmd.Flags().StringVarP(&createFilePath, "file", "f", "", "Path to a JSON file containing the portal definition")
	_ = CreateCmd.MarkFlagRequired("file")
}

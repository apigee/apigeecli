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

package apidocs

import (
	"fmt"
	"internal/apiclient"

	"internal/client/apidocs"

	"github.com/spf13/cobra"
)

// ImpCmd to import products
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import from a folder containing apidocs",
	Long:  "Import from a folder containing apidocs",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if siteid == "" {
			return fmt.Errorf("siteid is a mandatory parameter")
		}
		return apidocs.Import(siteid, folder)
	},
}

func init() {
	ImpCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "Folder containing site_<siteid>.json and apidocs_<siteid>_<id>.json files")

	_ = ImpCmd.MarkFlagRequired("file")
}

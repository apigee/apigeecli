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
	"fmt"

	"internal/apiclient"

	"internal/client/apps"

	"github.com/spf13/cobra"
)

// GetCmd to get app
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get App in an Organization by App ID",
	Long:  "Returns the app profile for the specified app ID.",
	Args: func(cmd *cobra.Command, args []string) error {
		if appID == "" && name == "" && productName == "" {
			return fmt.Errorf("pass either name or appId or productName")
		}
		if appID != "" && name != "" {
			return fmt.Errorf("name and appId cannot be used together")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if appID != "" {
			_, err = apps.Get(appID)
			return err
		}

		if name != "" {
			outBytes, err := apps.SearchApp(name)
			if err != nil {
				return err
			}
			// print the item
			return apiclient.PrettyPrint("json", outBytes)
		}

		_, err = apps.ListApps(productName)
		return err
	},
}

var productName string

func init() {
	GetCmd.Flags().StringVarP(&appID, "appId", "i",
		"", "Developer app id")
	GetCmd.Flags().StringVarP(&name, "name", "n",
		"", "Developer app name")
	GetCmd.Flags().StringVarP(&productName, "product", "p",
		"", "API Product name")
}

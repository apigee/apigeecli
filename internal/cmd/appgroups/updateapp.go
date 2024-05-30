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

package appgroups

import (
	"fmt"
	"strconv"

	"internal/apiclient"

	"internal/client/appgroups"

	"github.com/spf13/cobra"
)

// UpdateAppCmd to create app
var UpdateAppCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an App in an AppGroup",
	Long:  "Update an App in an AppGroup",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if expires != "" {
			if _, err = strconv.Atoi(expires); err != nil {
				return fmt.Errorf("expires must be an integer: %v", err)
			}
			expires += "000"
		}
		_, err = appgroups.UpdateApp(name, appName, expires, callback, apiProducts, scopes, attrs)
		return
	},
}

func init() {
	UpdateAppCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the AppGroup")
	UpdateAppCmd.Flags().StringVarP(&appName, "app-name", "",
		"", "Name of the app")
	UpdateAppCmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in seconds, for the lifetime of the consumer key")
	UpdateAppCmd.Flags().StringVarP(&callback, "callback", "c",
		"", "The callbackUrl is used by OAuth")
	UpdateAppCmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	UpdateAppCmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	UpdateAppCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = UpdateAppCmd.MarkFlagRequired("name")
	_ = UpdateAppCmd.MarkFlagRequired("app-name")
}

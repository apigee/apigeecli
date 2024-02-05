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

package apps

import (
	"fmt"
	"strconv"

	"internal/apiclient"

	"internal/client/apps"

	"github.com/spf13/cobra"
)

// UpdateCmd to update an app
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a Developer App",
	Long:  "Update a Developer App",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if expires != "" {
			if _, err = strconv.Atoi(expires); err != nil {
				return fmt.Errorf("expires must be an integer: %v", err)
			}
			expires += "000"
		}
		_, err = apps.Update(name, email, expires, callback, apiProducts, scopes, attrs)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	UpdateCmd.Flags().StringVarP(&email, "email", "e",
		"", "The developer's email")
	UpdateCmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in milliseconds, for the lifetime of the consumer key")
	UpdateCmd.Flags().StringVarP(&callback, "callback", "c",
		"", "The callbackUrl is used by OAuth")
	UpdateCmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	UpdateCmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	UpdateCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = UpdateCmd.MarkFlagRequired("name")
	_ = UpdateCmd.MarkFlagRequired("email")
}

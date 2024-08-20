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
	"strconv"

	"github.com/spf13/cobra"
)

// CreateCmd to create app
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Developer App",
	Long:  "Create a Developer App",
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
		_, err = apps.Create(name, email, expires, callback, apiProducts, scopes, attrs)
		return
	},
}

var (
	email, expires, callback string
	apiProducts, scopes      []string
	attrs                    map[string]string
)

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	CreateCmd.Flags().StringVarP(&email, "email", "e",
		"", "The developer's email or id")
	CreateCmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in seconds, for the lifetime of the consumer key")
	CreateCmd.Flags().StringVarP(&callback, "callback", "c",
		"", "The callbackUrl is used by OAuth")
	CreateCmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	CreateCmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	CreateCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("email")
}

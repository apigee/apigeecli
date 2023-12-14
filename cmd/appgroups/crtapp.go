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

// CreateAppCmd to create app
var CreateAppCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an App",
	Long:  "Create an App in an AppGroup",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if expires != "" {
			if _, err = strconv.Atoi(expires); err != nil {
				return fmt.Errorf("expires must be an integer: %v", err)
			}
			expires += "000"
		}
		_, err = appgroups.CreateApp(name, appName, expires, callback, apiProducts, scopes, attrs)
		return
	},
}

var (
	expires, callback   string
	apiProducts, scopes []string
)

func init() {
	CreateAppCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the AppGroup")
	CreateAppCmd.Flags().StringVarP(&appName, "app-name", "",
		"", "Name of the app")
	CreateAppCmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in seconds, for the lifetime of the consumer key")
	CreateAppCmd.Flags().StringVarP(&callback, "callback", "c",
		"", "The callbackUrl is used by OAuth")
	CreateAppCmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	CreateAppCmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	CreateAppCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = CreateAppCmd.MarkFlagRequired("name")
	_ = CreateAppCmd.MarkFlagRequired("app-name")
	_ = CreateAppCmd.MarkFlagRequired("prods")
}

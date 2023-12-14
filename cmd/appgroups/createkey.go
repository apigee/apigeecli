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

package appgroups

import (
	"fmt"
	"strconv"

	"internal/apiclient"
	"internal/client/appgroups"

	"github.com/spf13/cobra"
)

// CreateKeyCmd to create developer keys
var CreateKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an app key",
	Long:  "Create an app key",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if (key != "" && secret == "") || (secret != "" && key == "") {
			return fmt.Errorf("key and secret must both be passed or neither must be sent")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if expires != "" {
			if _, err = strconv.Atoi(expires); err != nil {
				return fmt.Errorf("expires must be an integer: %v", err)
			}
		}
		_, err = appgroups.CreateKey(name, appName, key, secret, expires, apiProducts, scopes, attrs)
		return
	},
}

var (
	secret  string
)

func init() {
	CreateKeyCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the AppGroup")
	CreateKeyCmd.Flags().StringVarP(&appName, "app-name", "",
		"", "Name of the app")
	CreateKeyCmd.Flags().StringVarP(&key, "key", "k",
		"", "Import an existing AppGroup app consumer key")
	CreateKeyCmd.Flags().StringVarP(&secret, "secret", "r",
		"", "Import an existing AppGroup app consumer secret")
	CreateKeyCmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in seconds, for the lifetime of the consumer key")
	CreateKeyCmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	CreateKeyCmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	CreateKeyCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = CreateKeyCmd.MarkFlagRequired("name")
	_ = CreateKeyCmd.MarkFlagRequired("app-name")
	_ = CreateKeyCmd.MarkFlagRequired("prods")
}

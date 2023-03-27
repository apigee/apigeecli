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
	"internal/apiclient"

	"internal/client/apps"

	"github.com/spf13/cobra"
)

// CreateKeyCmd to create developer keys
var CreateKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a developer app key",
	Long:  "Create a a developer app key",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apps.CreateKey(developerEmail, name, key, secret, apiProducts, scopes, attrs)
		return
	},
}

func init() {

	CreateKeyCmd.Flags().StringVarP(&key, "key", "k",
		"", "Developer app consumer key")
	CreateKeyCmd.Flags().StringVarP(&secret, "secret", "r",
		"", "Developer app consumer secret")
	CreateKeyCmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	CreateKeyCmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	CreateKeyCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = CreateKeyCmd.MarkFlagRequired("name")
	_ = CreateKeyCmd.MarkFlagRequired("key")
	_ = CreateKeyCmd.MarkFlagRequired("secret")
}

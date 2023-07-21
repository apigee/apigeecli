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
	"internal/apiclient"

	"internal/client/appgroups"

	"github.com/spf13/cobra"
)

// GetKeyCmd to get credential
var GetKeyCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets details for a consumer key for a AppGroup app",
	Long:  "Gets details for a consumer key for a AppGroup app",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = appgroups.GetKey(name, appName, key)
		return
	},
}

var key string

func init() {
	GetKeyCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the app group")
	GetKeyCmd.Flags().StringVarP(&appName, "app-name", "",
		"", "Name of the app")
	GetKeyCmd.Flags().StringVarP(&key, "key", "k",
		"", "App consumer key")

	_ = GetKeyCmd.MarkFlagRequired("name")
	_ = GetKeyCmd.MarkFlagRequired("app-name")
	_ = GetKeyCmd.MarkFlagRequired("key")
}

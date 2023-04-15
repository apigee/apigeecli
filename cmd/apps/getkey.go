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

// GetKeyCmd to manage tracing of apis
var GetKeyCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a developer app key",
	Long:  "Get a a developer app key",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apps.GetKey(developerEmail, name, key)
		return
	},
}

func init() {
	GetKeyCmd.Flags().StringVarP(&key, "key", "k",
		"", "developer app consumer key")

	_ = GetKeyCmd.MarkFlagRequired("key")
}

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

// CreateCmd to create appgroup
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an AppGroup",
	Long:  "Create an AppGroup",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = appgroups.Create(name, channelUri, channelId, displayName, attrs)
		return
	},
}

var channelId, channelUri, displayName string
var attrs map[string]string

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	CreateCmd.Flags().StringVarP(&channelId, "channelid", "i",
		"", "channel identifier identifies the owner maintaining this grouping")
	CreateCmd.Flags().StringVarP(&channelUri, "channelurl", "u",
		"", "A reference to the associated storefront/marketplace")
	CreateCmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "app group name displayed in the UI")
	CreateCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = CreateCmd.MarkFlagRequired("name")
}

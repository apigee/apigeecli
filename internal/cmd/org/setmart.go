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

package org

import (
	"internal/apiclient"

	"internal/client/orgs"

	"github.com/spf13/cobra"
)

// MartCmd to set mart endpoint
var MartCmd = &cobra.Command{
	Use:   "setmart",
	Short: "Set MART endpoint for an Apigee Org",
	Long:  "Set MART endpoint for an Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return orgs.SetOrgProperty("features.mart.server.endpoint", mart)
	},
}

var mart string

func init() {
	MartCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	MartCmd.Flags().StringVarP(&mart, "mart", "m",
		"", "MART Endpoint")

	_ = MartCmd.MarkFlagRequired("mart")
}

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
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/orgs"
)

//Cmd to get org details
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Apigee Org",
	Long:  "Create a new Apigee Org; Your GCP project must be whitelist for this operation",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		apiclient.SetProjectID(projectID)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = orgs.Create(region)
		return
	},
}

var region, projectID string

func init() {

	CreateCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	CreateCmd.Flags().StringVarP(&region, "reg", "r",
		"", "Analytics region name")
	CreateCmd.Flags().StringVarP(&projectID, "prj", "p",
		"", "GCP Project ID")

	_ = CreateCmd.MarkFlagRequired("prj")
	_ = CreateCmd.MarkFlagRequired("org")
	_ = CreateCmd.MarkFlagRequired("reg")
}

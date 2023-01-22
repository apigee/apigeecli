// Copyright 2022 Google LLC
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
	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/orgs"
	"github.com/spf13/cobra"
)

//DelCmd to get org details
var DelCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Apigee Org",
	Long:  "Delete an Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = orgs.Delete(retention)
		return
	},
}

var retention string

func init() {

	DelCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	DelCmd.Flags().StringVarP(&retention, "retension", "r",
		"", "Retention period for soft-delete; Must be MINIMUM or DELETION_RETENTION_UNSPECIFIED")
}

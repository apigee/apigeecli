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

// ListCmd to list orgs
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the Apigee organizations",
	Long:  "List the Apigee organizations, and the related projects that a user has permissions for",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = orgs.List()
		return
	},
}

func init() {
}

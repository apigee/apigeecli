// Copyright 2021 Google LLC
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

package eptattachment

import (
	"internal/apiclient"

	"internal/client/eptattachment"

	"github.com/spf13/cobra"
)

// ListCmd to list endpoint attachments
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List of service endpoints",
	Long:  "Lists of service endpoints",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = eptattachment.List()
		return
	},
}

// Copyright 2024 Google LLC
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

package projectregistrations

import (
	"internal/apiclient"
	"internal/client/hub"

	"github.com/spf13/cobra"
)

// CrtCmd
var CrtCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a host project registration",
	Long:  "Create a host project registration",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = hub.RegisterHostProject(registrationId, gcpProjectId)
		return
	},
}

var registrationId, gcpProjectId string

func init() {
	CrtCmd.Flags().StringVarP(&registrationId, "id", "i",
		"", "Registration ID")
	CrtCmd.Flags().StringVarP(&gcpProjectId, "project-id", "p",
		"", "Host Project ID")

	_ = CrtCmd.MarkFlagRequired("id")
	_ = CrtCmd.MarkFlagRequired("project-id")
}

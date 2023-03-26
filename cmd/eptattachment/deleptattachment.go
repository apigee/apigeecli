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

// Cmd to list endpoint attachments
var RemoveCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a service endpoint",
	Long:  "Delete a service endpoint",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = eptattachment.Delete(name)
		return
	},
}

func init() {

	RemoveCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the service endpoint")

	_ = RemoveCmd.MarkFlagRequired("name")
}

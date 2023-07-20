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

package developers

import (
	"internal/apiclient"

	"internal/client/developers"

	"github.com/spf13/cobra"
)

// UpdateCmd to update developer
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an Apigee developer configuration",
	Long:  "Update an Apigee developer configuration",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = developers.Update(email, firstName, lastName, userName, status, attrs)
		return
	},
}

var status bool

func init() {
	UpdateCmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")
	UpdateCmd.Flags().StringVarP(&firstName, "first", "f",
		"", "The first name of the developer")
	UpdateCmd.Flags().StringVarP(&lastName, "last", "s",
		"", "The last name of the developer")
	UpdateCmd.Flags().StringVarP(&userName, "user", "u",
		"", "The username of the developer")
	UpdateCmd.Flags().BoolVarP(&status, "status", "",
		true, "The status of the developer")
	UpdateCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = UpdateCmd.MarkFlagRequired("email")
}

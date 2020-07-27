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
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/developers"
)

//Cmd to create developer
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a developer",
	Long:  "Create a developer",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = developers.Create(email, firstName, lastName, userName, attrs)
		return
	},
}

var lastName, firstName, userName string
var attrs map[string]string

func init() {

	CreateCmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")
	CreateCmd.Flags().StringVarP(&firstName, "first", "f",
		"", "The first name of the developer")
	CreateCmd.Flags().StringVarP(&lastName, "last", "s",
		"", "The last name of the developer")
	CreateCmd.Flags().StringVarP(&userName, "user", "u",
		"", "The username of the developer")
	CreateCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = CreateCmd.MarkFlagRequired("email")
	_ = CreateCmd.MarkFlagRequired("first")
	_ = CreateCmd.MarkFlagRequired("last")
	_ = CreateCmd.MarkFlagRequired("user")
}

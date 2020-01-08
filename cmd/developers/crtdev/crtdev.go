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

package crtdev

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/client/developers"
)

//Cmd to create developer
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a developer",
	Long:  "Create a developer",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = developers.Create(email, firstName, lastName, userName, attrs)
		return
	},
}

var email, lastName, firstName, userName string
var attrs map[string]string

func init() {

	Cmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")
	Cmd.Flags().StringVarP(&firstName, "first", "f",
		"", "The first name of the developer")
	Cmd.Flags().StringVarP(&lastName, "last", "s",
		"", "The last name of the developer")
	Cmd.Flags().StringVarP(&userName, "user", "u",
		"", "The username of the developer")
	Cmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = Cmd.MarkFlagRequired("email")
	_ = Cmd.MarkFlagRequired("first")
	_ = Cmd.MarkFlagRequired("last")
	_ = Cmd.MarkFlagRequired("user")
}

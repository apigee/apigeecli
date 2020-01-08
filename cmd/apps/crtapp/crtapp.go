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

package crtapp

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/client/apps"
)

//Cmd to create app
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Developer App",
	Long:  "Create a Developer App",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apps.Create(name, email, expires, callback, apiProducts, scopes, attrs)
		return
	},
}

var name, email, expires, callback string
var apiProducts, scopes []string
var attrs map[string]string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	Cmd.Flags().StringVarP(&email, "email", "e",
		"", "The developer's email or id")
	Cmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in milliseconds, for the lifetime of the consumer key")
	Cmd.Flags().StringVarP(&callback, "callback", "c",
		"", "The callbackUrl is used by OAuth")
	Cmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	Cmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	Cmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("email")
}

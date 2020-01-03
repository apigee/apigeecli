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

package enableac

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/org/setprop"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to set mart endpoint
var Cmd = &cobra.Command{
	Use:   "enable-apigee-connect",
	Short: "Set MART endpoint for an Apigee Org",
	Long:  "Set MART endpoint for an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return setprop.SetOrgProperty("features.mart.apigee.connect.enabled", "true")
	},
}

var mart string
var whitelist bool

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkFlagRequired("org")
}

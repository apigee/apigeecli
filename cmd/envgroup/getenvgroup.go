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

package envgroup

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/envgroups"
)

//GetCmd to get env group
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets an Environment Group",
	Long:  "Gets an Environment Group",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = envgroups.Get(name)
		return
	},
}

func init() {

	GetCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	GetCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the environment group")

	_ = GetCmd.MarkFlagRequired("org")
	_ = GetCmd.MarkFlagRequired("name")
}

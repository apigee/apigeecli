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
	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/envgroups"
	"github.com/spf13/cobra"
)

//Cmd to create a new product
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Environment Group",
	Long:  "Create an Environment Group",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = envgroups.Create(name, hostnames)
		return
	},
}

func init() {

	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Environment Group")
	CreateCmd.Flags().StringArrayVarP(&hostnames, "hosts", "d",
		[]string{}, "A list of hostnames")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("hosts")
}

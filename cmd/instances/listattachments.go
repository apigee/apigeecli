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

package instances

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/instances"
)

//Cmd to create a new instance
var ListAttachCmd = &cobra.Command{
	Use:   "list",
	Short: "List attachment details for an instance",
	Long:  "List attachment details for an instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = instances.ListAttach(name)
		return
	},
}

func init() {

	ListAttachCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Instance")

	_ = ListAttachCmd.MarkFlagRequired("name")
}

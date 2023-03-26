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

package ops

import (
	"internal/apiclient"

	"internal/client/operations"

	"github.com/spf13/cobra"
)

// Cmd to get ops details
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get operation of an org",
	Long:  "Get operation of an org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = operations.Get(name)
		return
	},
}

var name string

func init() {

	GetCmd.Flags().StringVarP(&name, "name", "n",
		"", "operation name")

	_ = GetCmd.MarkFlagRequired("name")
}

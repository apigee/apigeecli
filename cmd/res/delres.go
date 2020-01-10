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

package res

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/res"
)

//Cmd to del a resource
var DelCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource file",
	Long:  "Delete a resource file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = res.Delete(name, resType)
		return
	},
}

func init() {

	DelCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the resource file")
	DelCmd.Flags().StringVarP(&resType, "type", "p",
		"", "Resource type")

	_ = DelCmd.MarkFlagRequired("name")
	_ = DelCmd.MarkFlagRequired("type")
}

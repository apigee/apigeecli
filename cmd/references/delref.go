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

package references

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/references"
)

//Cmd to delete reference
var DelCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a reference in an environment",
	Long:  "Delete a reference in an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetApigeeOrg(org)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = references.Delete(name)
		return
	},
}

func init() {

	DelCmd.PersistentFlags().StringVarP(&name, "name", "n",
		"", "Reference name")

	_ = DelCmd.MarkFlagRequired("name")
}

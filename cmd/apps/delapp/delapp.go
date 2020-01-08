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

package deldev

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/client/apps"
)

//Cmd to delete app
var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a Developer App from an organization",
	Long:  "Deletes a Developer Appfrom an organization",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apps.Delete(name, id)
		return
	},
}

var name, id string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	Cmd.Flags().StringVarP(&id, "id", "i",
		"", "Developer Id")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("id")
}

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

package apps

import (
	"github.com/spf13/cobra"
)

//KeysCmd
var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage developer app keys",
	Long:  "Manage developer apps keys",
}

var developerEmail, key string

func init() {

	KeysCmd.PersistentFlags().StringVarP(&name, "name", "n",
		"", "Developer App name")

	KeysCmd.PersistentFlags().StringVarP(&developerEmail, "dev", "d",
		"", "Developer email")

	_ = KeysCmd.MarkPersistentFlagRequired("dev")
	_ = KeysCmd.MarkPersistentFlagRequired("app")

	KeysCmd.AddCommand(CreateKeyCmd)
	KeysCmd.AddCommand(GetKeyCmd)
	KeysCmd.AddCommand(DeleteKeyCmd)
	KeysCmd.AddCommand(UpdateKeyCmd)
}

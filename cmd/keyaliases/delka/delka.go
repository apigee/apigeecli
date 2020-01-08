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

package delka

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/client/keyaliases"
)

//Cmd to delete key aliases
var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Key Alias",
	Long:  "Delete a Key Alias",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = keyaliases.Delete(keystoreName, aliasName)
		return
	},
}

var keystoreName, aliasName string

func init() {
	Cmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")
	Cmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key alias")

	_ = Cmd.MarkFlagRequired("alias")
	_ = Cmd.MarkFlagRequired("key")
}

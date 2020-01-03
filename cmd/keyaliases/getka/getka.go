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

package getka

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to get key aliases
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get a Key Alias",
	Long:  "Get a Key Alias",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env,
			"keystores", shared.RootArgs.AliasName, "aliases", aliasName)
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var aliasName string

func init() {

	Cmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key store")

	_ = Cmd.MarkFlagRequired("alias")
}

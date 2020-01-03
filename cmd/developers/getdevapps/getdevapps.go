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

package getdevapps

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

//Cmd to get developer
var Cmd = &cobra.Command{
	Use:   "getapps",
	Short: "Returns the apps owned by a developer by email address",
	Long:  "Returns the apps owned by a developer by email address",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		if expand {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", name, "apps")
			q := u.Query()
			q.Set("expand", "true")
			u.RawQuery = q.Encode()
		} else {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", name, "apps")
		}
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var name string
var expand bool

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "email of the developer")

	Cmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "expand app details")

	_ = Cmd.MarkFlagRequired("org")
	_ = Cmd.MarkFlagRequired("name")
}

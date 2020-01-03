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

package listdev

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to list developer
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of App Developers",
	Long:  "Lists all developers in an organization by email address",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers")
		q := u.Query()
		if expand {
			q.Set("expand", "true")
		} else {
			q.Set("expand", "false")
		}
		if count != "" {
			q.Set("count", count)
		}
		u.RawQuery = q.Encode()
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var expand = false
var count string

func init() {

	Cmd.Flags().StringVarP(&count, "count", "c",
		"", "Number of developers; limit is 1000")

	Cmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "Expand Details")
}

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

package gettrcapi

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get a debug session for an API proxy revision",
	Long:  "Get a debug session for an API proxy revision deployed in an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		if messageID == "" {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env,
				"apis", name, "revisions", revision, "debugsessions", sessionID, "data")
			q := u.Query()
			q.Set("limit", "20")
			u.RawQuery = q.Encode()
		} else {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env,
				"apis", name, "revisions", revision, "debugsessions", sessionID, "data", messageID)
		}

		_, err = shared.HttpClient(true, u.String())
		return

	},
}

var name, revision, sessionID, messageID string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "API Proxy revision")
	Cmd.Flags().StringVarP(&sessionID, "ses", "s",
		"", "Debug session Id")
	Cmd.Flags().StringVarP(&messageID, "msg", "m",
		"", "Debug session Id")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("rev")
	_ = Cmd.MarkFlagRequired("ses")
}

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

package depsf

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to deploy shared flow
var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a revision of an existing Sharedflow",
	Long:  "Deploys a revision of an existing Sharedflow to an environment in an organization",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		if overrides {
			q := u.Query()
			q.Set("override", "true")
			u.RawQuery = q.Encode()
		}
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env,
			"sharedflows", name, "revisions", revision, "deployments")
		_, err = shared.HttpClient(true, u.String(), "")
		return
	},
}

var name, revision string
var overrides bool

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Sharedflow name")
	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "Sharedflow revision")
	Cmd.Flags().BoolVarP(&overrides, "ovr", "r",
		false, "Forces deployment of the new revision")

	_ = Cmd.MarkFlagRequired("env")
	_ = Cmd.MarkFlagRequired("name")
}

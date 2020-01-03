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

package listres

import (
	"fmt"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/srinandan/apigeecli/cmd/types"
)

//Cmd to list resources
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List all resources in your environment",
	Long:  "List all resources in your environment",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {

		if resType != "" && !types.IsValidResource(resType) {
			return fmt.Errorf("invalid resource type")
		}
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "resourcefiles")

		if resType != "" {
			q := u.Query()
			q.Set("type", resType)
			u.RawQuery = q.Encode()
		}
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var resType string

func init() {

	Cmd.Flags().StringVarP(&resType, "type", "p",
		"", "Resource type")

}

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

package setdebugmask

import (
	"encoding/json"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "set",
	Short: "Set debugmasks for an Environment",
	Long:  "Set debugmasks for an Environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		//the following steps will validate json
		m := map[string]string{}
		err = json.Unmarshal([]byte(payload), &m)
		if err != nil {
			return err
		}
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "debugmask")
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}

var payload string

func init() {

	Cmd.Flags().StringVarP(&payload, "mask", "m",
		"", "Mask configuration is in JSON format")

	_ = Cmd.MarkFlagRequired("mask")
}

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

package getdev

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/thedevsaddam/gojsonq"
)

//Cmd to get app
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get App in an Organization by App ID",
	Long:  "Returns the app profile for the specified app ID.",
	Args: func(cmd *cobra.Command, args []string) error {
		if appID == "" && name == "" {
			return fmt.Errorf("pass either name or appId")
		}
		if appID != "" && name != "" {
			return fmt.Errorf("name and appId cannot be used together")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		if appID != "" {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "apps", appID)
			_, err = shared.HttpClient(true, u.String())
			return
		}
		//search by name is not implement; use list and return the appropriate app
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apps")
		q := u.Query()
		q.Set("expand", "true")
		q.Set("includeCred", "false")
		u.RawQuery = q.Encode()
		//don't print the list
		respBody, err := shared.HttpClient(false, u.String())
		if err != nil {
			return err
		}
		jq := gojsonq.New().JSONString(string(respBody)).From("app").Where("name", "eq", name)
		out := jq.Get()
		outBytes, err := json.Marshal(out)
		if err != nil {
			return err
		}
		//print the item
		return shared.PrettyPrint(outBytes)
	},
}

var appID, name string

func init() {

	Cmd.Flags().StringVarP(&appID, "appId", "i",
		"", "Developer app id")
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Developer app name")
}

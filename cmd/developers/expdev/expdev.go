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

package expdev

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to export developer
var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export Developers to a file",
	Long:  "Export Developers to a file",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		const exportFileName = "developers.json"

		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers")

		q := u.Query()
		q.Set("expand", "true")

		u.RawQuery = q.Encode()
		//don't print to sysout
		respBody, err := shared.HttpClient(false, u.String())
		if err != nil {
			return err
		}

		return shared.WriteByteArrayToFile(exportFileName, false, respBody)
	},
}

func init() {

}

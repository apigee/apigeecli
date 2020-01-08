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
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/developers"
)

//Cmd to export developer
var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export Developers to a file",
	Long:  "Export Developers to a file",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		const exportFileName = "developers.json"

		respBody, err := developers.Export()
		if err != nil {
			return err
		}

		return apiclient.WriteByteArrayToFile(exportFileName, false, respBody)
	},
}

func init() {

}

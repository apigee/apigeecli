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

package targetservers

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/targetservers"
)

//ExpCmd to export target servers
var ExpCmd = &cobra.Command{
	Use:   "export",
	Short: "Export target servers to a file",
	Long:  "Export target servers to a file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		const exportFileName = "targetservers.json"
		payload, err := targetservers.Export(conn)
		if err != nil {
			return
		}
		return apiclient.WriteArrayByteArrayToFile(exportFileName, false, payload)
	},
}

var conn int

func init() {

	ExpCmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

}

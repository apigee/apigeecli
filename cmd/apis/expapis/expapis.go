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

package expapis

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/client/apis"
)

//Cmd to export apis
var Cmd = &cobra.Command{
	Use:   "export",
	Short: "export API proxy bundles from an org",
	Long:  "export API proxy bundles from an org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return apis.ExportProxies(conn)
	},
}

var conn int

func init() {

	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
}

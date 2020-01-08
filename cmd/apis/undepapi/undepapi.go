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

package undepapi

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/apis"
)

//Cmd to undeloy api
var Cmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploys a revision of an existing API proxy",
	Long:  "Undeploys a revision of an existing API proxy to an environment in an organization",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.UndeployProxy(name, revision)
		return
	},
}

var name string
var revision int

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	Cmd.Flags().StringVarP(apiclient.GetApigeeEnvP(), "env", "e",
		"", "Apigee environment name")
	Cmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision")

	_ = Cmd.MarkFlagRequired("env")
	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("rev")

}

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
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/sharedflows"
)

//Cmd to deploy shared flow
var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a revision of an existing Sharedflow",
	Long:  "Deploys a revision of an existing Sharedflow to an environment in an organization",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = sharedflows.Deploy(name, revision, overrides)
		return
	},
}

var name string
var revision int
var overrides bool

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Sharedflow name")
	Cmd.Flags().StringVarP(apiclient.GetApigeeEnvP(), "env", "e",
		"", "Apigee environment name")
	Cmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "Sharedflow revision")
	Cmd.Flags().BoolVarP(&overrides, "ovr", "r",
		false, "Forces deployment of the new revision")

	_ = Cmd.MarkFlagRequired("env")
	_ = Cmd.MarkFlagRequired("name")
}

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

package getenv

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/env"
)

//Cmd to get env details
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get properties of an environment",
	Long:  "Get properties of an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = env.Get(config)
		return
	},
}

var config = false

func init() {

	Cmd.Flags().StringVarP(apiclient.GetApigeeEnvP(), "env", "e",
		"", "Apigee environment name")
	Cmd.Flags().BoolVarP(&config, "config", "c",
		false, "Return configuration details")

	_ = Cmd.MarkFlagRequired("env")
}

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

package iam

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/env/iam/getiam"
	"github.com/srinandan/apigeecli/cmd/env/iam/setax"
	"github.com/srinandan/apigeecli/cmd/env/iam/setdeploy"
	"github.com/srinandan/apigeecli/cmd/env/iam/setsync"
	"github.com/srinandan/apigeecli/cmd/env/iam/testiam"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM permissions for the environment",
	Long:  "Manage IAM permissions for the environment",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(getiam.Cmd)
	Cmd.AddCommand(setax.Cmd)
	Cmd.AddCommand(setdeploy.Cmd)
	Cmd.AddCommand(setsync.Cmd)
	Cmd.AddCommand(testiam.Cmd)
}

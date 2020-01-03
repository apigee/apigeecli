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

package env

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/env/debugmask"
	"github.com/srinandan/apigeecli/cmd/env/getenv"
	"github.com/srinandan/apigeecli/cmd/env/iam"
	"github.com/srinandan/apigeecli/cmd/env/listenv"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage envs
var Cmd = &cobra.Command{
	Use:     "envs",
	Aliases: []string{"environments"},
	Short:   "Manage Apigee environments",
	Long:    "Manage Apigee environments",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listenv.Cmd)
	Cmd.AddCommand(getenv.Cmd)
	Cmd.AddCommand(iam.Cmd)
	Cmd.AddCommand(debugmask.Cmd)
}

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
	"github.com/srinandan/apigeecli/cmd/flowhooks/crtfh"
	"github.com/srinandan/apigeecli/cmd/flowhooks/delfh"
	"github.com/srinandan/apigeecli/cmd/flowhooks/getfh"
	"github.com/srinandan/apigeecli/cmd/flowhooks/listfh"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage flow hooks
var Cmd = &cobra.Command{
	Use:     "flowhooks",
	Aliases: []string{"fs"},
	Short:   "Manage Flowhooks",
	Long:    "Manage Flowhooks",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(listfh.Cmd)
	Cmd.AddCommand(getfh.Cmd)
	Cmd.AddCommand(delfh.Cmd)
	Cmd.AddCommand(crtfh.Cmd)
}

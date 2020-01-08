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

package res

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/cmd/res/crtres"
	"github.com/srinandan/apigeecli/cmd/res/delres"
	"github.com/srinandan/apigeecli/cmd/res/getres"
	"github.com/srinandan/apigeecli/cmd/res/listres"
)

//Cmd to manage resources
var Cmd = &cobra.Command{
	Use:     "resources",
	Aliases: []string{"res"},
	Short:   "Manage resources within an Apigee environment",
	Long:    "Manage resources within an Apigee environment",
}

func init() {

	Cmd.PersistentFlags().StringVarP(apiclient.GetApigeeOrgP(), "org", "o",
		"", "Apigee organization name")

	Cmd.PersistentFlags().StringVarP(apiclient.GetApigeeEnvP(), "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(crtres.Cmd)
	Cmd.AddCommand(delres.Cmd)
	Cmd.AddCommand(listres.Cmd)
	Cmd.AddCommand(getres.Cmd)
}

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
	"github.com/srinandan/apigeecli/cmd/targetservers/crtts"
	"github.com/srinandan/apigeecli/cmd/targetservers/delts"
	"github.com/srinandan/apigeecli/cmd/targetservers/expts"
	"github.com/srinandan/apigeecli/cmd/targetservers/getts"
	"github.com/srinandan/apigeecli/cmd/targetservers/listts"
)

//Cmd to manage targetservers
var Cmd = &cobra.Command{
	Use:     "targetservers",
	Aliases: []string{"ts"},
	Short:   "Manage Target Servers",
	Long:    "Manage Target Servers",
}

func init() {

	Cmd.PersistentFlags().StringVarP(apiclient.GetApigeeOrgP(), "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(apiclient.GetApigeeEnvP(), "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(listts.Cmd)
	Cmd.AddCommand(getts.Cmd)
	Cmd.AddCommand(delts.Cmd)
	Cmd.AddCommand(expts.Cmd)
	Cmd.AddCommand(crtts.Cmd)
}

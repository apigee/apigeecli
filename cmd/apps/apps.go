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

package apps

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	crtapp "github.com/srinandan/apigeecli/cmd/apps/crtapp"
	delapp "github.com/srinandan/apigeecli/cmd/apps/delapp"
	expapp "github.com/srinandan/apigeecli/cmd/apps/expapp"
	genkey "github.com/srinandan/apigeecli/cmd/apps/genkey"
	getapp "github.com/srinandan/apigeecli/cmd/apps/getapp"
	impapps "github.com/srinandan/apigeecli/cmd/apps/impapps"
	listapp "github.com/srinandan/apigeecli/cmd/apps/listapp"
)

//Cmd to manage apps
var Cmd = &cobra.Command{
	Use:     "apps",
	Aliases: []string{"applications"},
	Short:   "Manage Apigee Developer Applications",
	Long:    "Manage Apigee Developer Applications",
}

func init() {

	Cmd.PersistentFlags().StringVarP(apiclient.GetApigeeOrgP(), "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listapp.Cmd)
	Cmd.AddCommand(getapp.Cmd)
	Cmd.AddCommand(delapp.Cmd)
	Cmd.AddCommand(crtapp.Cmd)
	Cmd.AddCommand(genkey.Cmd)
	Cmd.AddCommand(impapps.Cmd)
	Cmd.AddCommand(expapp.Cmd)
}

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

package developers

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/developers/crtdev"
	"github.com/srinandan/apigeecli/cmd/developers/deldev"
	"github.com/srinandan/apigeecli/cmd/developers/expdev"
	"github.com/srinandan/apigeecli/cmd/developers/getdev"
	"github.com/srinandan/apigeecli/cmd/developers/getdevapps"
	"github.com/srinandan/apigeecli/cmd/developers/impdev"
	"github.com/srinandan/apigeecli/cmd/developers/listdev"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage developers
var Cmd = &cobra.Command{
	Use:     "devs",
	Aliases: []string{"developers"},
	Short:   "Manage Apigee App Developers",
	Long:    "Manage Apigee App Developers",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listdev.Cmd)
	Cmd.AddCommand(getdev.Cmd)
	Cmd.AddCommand(deldev.Cmd)
	Cmd.AddCommand(crtdev.Cmd)
	Cmd.AddCommand(impdev.Cmd)
	Cmd.AddCommand(expdev.Cmd)
	Cmd.AddCommand(getdevapps.Cmd)
}

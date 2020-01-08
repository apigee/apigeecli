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

package sharedflows

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/cmd/sharedflows/crtsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/delsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/depsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/expsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/fetchsf"
	getsf "github.com/srinandan/apigeecli/cmd/sharedflows/getsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/impsfs"
	listsf "github.com/srinandan/apigeecli/cmd/sharedflows/listsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/undepsf"
)

//Cmd to manage shared flows
var Cmd = &cobra.Command{
	Use:   "sharedflows",
	Short: "Manage Apigee shared flows in an org",
	Long:  "Manage Apigee shared flows in an org",
}

func init() {

	Cmd.PersistentFlags().StringVarP(apiclient.GetApigeeOrgP(), "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listsf.Cmd)
	Cmd.AddCommand(getsf.Cmd)
	Cmd.AddCommand(fetchsf.Cmd)
	Cmd.AddCommand(delsf.Cmd)
	Cmd.AddCommand(undepsf.Cmd)
	Cmd.AddCommand(depsf.Cmd)
	Cmd.AddCommand(crtsf.Cmd)
	Cmd.AddCommand(expsf.Cmd)
	Cmd.AddCommand(impsfs.Cmd)
}

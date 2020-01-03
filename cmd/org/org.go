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

package org

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/org/createorg"
	"github.com/srinandan/apigeecli/cmd/org/enableac"
	"github.com/srinandan/apigeecli/cmd/org/enablewl"
	"github.com/srinandan/apigeecli/cmd/org/getorg"
	"github.com/srinandan/apigeecli/cmd/org/listorgs"
	"github.com/srinandan/apigeecli/cmd/org/setmart"
)

//Cmd to manage orgs
var Cmd = &cobra.Command{
	Use:   "orgs",
	Short: "Manage Apigee Orgs",
	Long:  "Manage Apigee Orgs",
}

func init() {
	Cmd.AddCommand(createorg.Cmd)
	Cmd.AddCommand(listorgs.Cmd)
	Cmd.AddCommand(getorg.Cmd)
	Cmd.AddCommand(setmart.Cmd)
	Cmd.AddCommand(enableac.Cmd)
	Cmd.AddCommand(enablewl.Cmd)
}

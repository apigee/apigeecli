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

package cache

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/cmd/cache/delcache"
	"github.com/srinandan/apigeecli/cmd/cache/listcache"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage caches within an Apigee environment",
	Long:  "Manage caches within an Apigee environment",
}

func init() {

	Cmd.PersistentFlags().StringVarP(apiclient.GetApigeeOrgP(), "org", "o",
		"", "Apigee organization name")

	Cmd.PersistentFlags().StringVarP(apiclient.GetApigeeEnvP(), "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(listcache.Cmd)
	Cmd.AddCommand(delcache.Cmd)
}

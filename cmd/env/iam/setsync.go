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
	"github.com/srinandan/apigeecli/apiclient"
	environments "github.com/srinandan/apigeecli/client/env"
)

//Cmd to manage tracing of apis
var SetSyncCmd = &cobra.Command{
	Use:   "setsync",
	Short: "Set Synchronization Manager role for a SA on an environment",
	Long:  "Set Synchronization Manager role for a SA on an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return environments.SetIAM(serviceAccountName, "sync")
	},
}

func init() {

	SetSyncCmd.Flags().StringVarP(&serviceAccountName, "name", "n",
		"", "Service Account Name")

	_ = Cmd.MarkFlagRequired("name")
}

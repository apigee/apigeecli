// Copyright 2021 Google LLC
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
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/env"
)

//CreateArchiveCmd to create env archive
var CreateArchiveCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new revision of archive in the environment",
	Long:  "Create a new revision of archive in the environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = env.CreateArchive(name, zipfile)
		return
	},
}

var zipfile string

func init() {
	CreateArchiveCmd.Flags().StringVarP(&name, "name", "n",
		"", "Archive name")
	CreateArchiveCmd.Flags().StringVarP(&zipfile, "zipfile", "z",
		"", "Archive Zip file")

	_ = CreateArchiveCmd.MarkFlagRequired("name")
	_ = CreateArchiveCmd.MarkFlagRequired("zipfile")
}

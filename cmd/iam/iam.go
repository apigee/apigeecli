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
)

//Cmd to manage orgs
var Cmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM permissions for Apigee",
	Long: "Manage IAM permissions for Apigee. The SA to run this command requires Security Admin, " +
		"Create Service Accounts and Service Account Key Admin roles",
}

var name, projectID, roleType string
var generateName, wid bool

var roles = []string{"mart", "analytics", "all", "logger", "connect", "cassandra", "watcher", "sync"}

func init() {
	Cmd.AddCommand(CallCmd)
}

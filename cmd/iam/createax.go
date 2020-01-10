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
)

//Cmd to get org details
var CaxCmd = &cobra.Command{
	Use:   "createax",
	Short: "Create a new IAM Service Account for Apigee Analytics",
	Long:  "Create a new IAM Service Account for Apigee Analytics",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(projectID)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.CreateIAMServiceAccount(name, "analytics")
	},
}

func init() {

	CaxCmd.Flags().StringVarP(&projectID, "prj", "p",
		"", "GCP Project ID")
	CaxCmd.Flags().StringVarP(&name, "name", "n",
		"", "Service Account Name")

	_ = CaxCmd.MarkFlagRequired("prj")
	_ = CaxCmd.MarkFlagRequired("name")
}

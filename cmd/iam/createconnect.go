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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
)

//Cmd to get org details
var ConnCmd = &cobra.Command{
	Use:   "createconnect",
	Short: "Create a new IAM Service Account for Apigee Connect",
	Long:  "Create a new IAM Service Account for Apigee Connect",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if !generateName && name == "" {
			return fmt.Errorf("provide a service account name or allow the tool to generate one")
		}
		apiclient.SetProjectID(projectID)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if generateName {
			name = GenerateName("apigee-conn-")
		}
		return apiclient.CreateIAMServiceAccount(name, "connect")
	},
}

func init() {

	ConnCmd.Flags().StringVarP(&projectID, "prj", "p",
		"", "GCP Project ID")
	ConnCmd.Flags().StringVarP(&name, "name", "n",
		"", "Service Account Name")
	ConnCmd.Flags().BoolVarP(&generateName, "gen", "g",
		false, "Generate account name")

	_ = ConnCmd.MarkFlagRequired("prj")
}

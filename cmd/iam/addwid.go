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
	"github.com/apigee/apigeecli/apiclient"
	"github.com/spf13/cobra"
)

//WidCmd to get org details
var WidCmd = &cobra.Command{
	Use:   "addwid",
	Short: "Provide WID role to a Service Account for Apigee Runtime",
	Long:  "Provide Workload Identity user role to an IAM Service Account for Apigee Runtime",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(projectID)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.AddWid(projectID, namespace, kServiceAccount, gServiceAccount)
	},
}

var namespace, kServiceAccount, gServiceAccount string

func init() {

	WidCmd.Flags().StringVarP(&projectID, "prj", "p",
		"", "GCP Project ID")
	WidCmd.Flags().StringVarP(&namespace, "namespace", "n",
		"apigee", "Apigee runtime namespace")
	WidCmd.Flags().StringVarP(&gServiceAccount, "gsa", "g",
		"", "Google Service Account Name")
	WidCmd.Flags().StringVarP(&kServiceAccount, "ksa", "k",
		"", "Kubernetes Service Account Name")

	_ = WidCmd.MarkFlagRequired("prj")
	_ = WidCmd.MarkFlagRequired("namespace")
	_ = WidCmd.MarkFlagRequired("gsa")
	_ = WidCmd.MarkFlagRequired("ksa")
}

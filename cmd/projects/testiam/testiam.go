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

package testiam

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "testiam",
	Short: "Test IAM policy for a GCP Project",
	Long:  "Test IAM policy for a GCP Project",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(apiclient.CrmURL)
		u.Path = path.Join(u.Path, apiclient.GetProjectID()+":testIamPermissions")
		payload := "{\"permissions\":[\"resourcemanager.projects.get\"]}"
		_, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
		return
	},
}

func init() {

	Cmd.Flags().StringVarP(apiclient.GetProjectIDP(), "prj", "p",
		"", "GCP Project ID")

	_ = Cmd.MarkFlagRequired("prj")
}

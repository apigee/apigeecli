// Copyright 2022 Google LLC
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

package apis

import (
	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/apis"
	"github.com/spf13/cobra"
)

//UpdateCmd to list api
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update APIs in an Apigee Org",
	Long:  "Update APIs in an Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apis.Update(name, labels)
		return
	},
}

var labels map[string]string

func init() {
	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	UpdateCmd.Flags().StringToStringVar(&labels, "labels",
		nil, "Labels")

	_ = UpdateCmd.MarkFlagRequired("name")
	_ = UpdateCmd.MarkFlagRequired("labels")
}

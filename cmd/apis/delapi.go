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

package apis

import (
	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/apis"
	"github.com/spf13/cobra"
)

//DelCmd to delete api
var DelCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an API proxy",
	Long: "Deletes an API proxy and all associated endpoints, policies, resources, and revisions." +
		"The proxy must be undeployed first.",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if revision == -1 {
			_, err = apis.DeleteProxy(name)
		} else {
			_, err = apis.DeleteProxyRevision(name, revision)
		}
		return
	},
}

func init() {
	DelCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	DelCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision")
	_ = DelCmd.MarkFlagRequired("name")
}

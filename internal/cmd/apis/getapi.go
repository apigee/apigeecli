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
	"fmt"
	"internal/apiclient"
	"internal/client/apis"

	"github.com/spf13/cobra"
)

// GetCmd to get api details
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets an API Proxy by name",
	Long:  "Gets an API Proxy by name, including a list of its revisions.",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if name == "" {
			return fmt.Errorf("name cannot be empty")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if name != "" {
			_, err = apis.GetProxy(name, revision)
			return err
		}
		return
	},
}

func init() {
	GetCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	GetCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision")

	_ = GetCmd.MarkFlagRequired("name")
}

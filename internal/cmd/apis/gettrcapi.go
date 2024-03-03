// Copyright 2020-2024 Google LLC
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
	"internal/apiclient"

	"internal/client/apis"

	"github.com/spf13/cobra"
)

// GetTrcCmd to manage tracing of apis
var GetTrcCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a debug session for an API proxy revision",
	Long:  "Get a debug session for an API proxy revision deployed in an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if revision == -1 {
			if revision, err = apis.GetHighestProxyRevision(name); err != nil {
				return err
			}
		}
		_, err = apis.GetTraceSession(name, revision, sessionID, messageID)
		return
	},
}

var sessionID, messageID string

func init() {
	GetTrcCmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	GetTrcCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "API Proxy revision. If not set, the highest revision is used")
	GetTrcCmd.Flags().StringVarP(&sessionID, "ses", "s",
		"", "Debug session Id")
	GetTrcCmd.Flags().StringVarP(&messageID, "msg", "m",
		"", "Debug transaction Id")

	_ = GetTrcCmd.MarkFlagRequired("name")
	_ = GetTrcCmd.MarkFlagRequired("ses")
}

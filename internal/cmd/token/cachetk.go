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

package token

import (
	"internal/apiclient"

	"internal/clilog"

	"github.com/spf13/cobra"
)

// CacheCmd to cache token
var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Generate and cache a new access token",
	Long:  "Generate and cache a new access token",
	Args: func(cmd *cobra.Command, args []string) error {
		apiclient.SetServiceAccount(serviceAccount)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		clilog.Init(apiclient.DebugEnabled(), apiclient.GetPrintOutput(), apiclient.GetNoOutput())
		err := apiclient.SetAccessToken()
		clilog.Info.Printf("Token %s cached\n", apiclient.GetApigeeToken())
		return err
	},
}

func init() {
	_ = CacheCmd.MarkFlagRequired("account")
}

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

package apps

import (
	"fmt"
	"strconv"

	"internal/apiclient"

	"internal/client/apps"

	"github.com/spf13/cobra"
)

// GenKeyCmd to generate key
var GenKeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generate a new developer KeyPair",
	Long:  "Generate a new developer KeyPair",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if expires != "" {
			if _, err = strconv.Atoi(expires); err != nil {
				return fmt.Errorf("expires must be an integer: %v", err)
			}
			expires += "000"
		}
		_, err = apps.GenerateKey(name, devID, apiProducts, callback, expires, scopes)
		return
	},
}

var devID string

func init() {
	GenKeyCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	GenKeyCmd.Flags().StringVarP(&devID, "devid", "d",
		"", "The developer's id or email")
	GenKeyCmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in seconds, for the lifetime of the consumer key")
	GenKeyCmd.Flags().StringVarP(&callback, "callback", "c",
		"", "The callbackUrl is used by OAuth")
	GenKeyCmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	GenKeyCmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")

	_ = GenKeyCmd.MarkFlagRequired("name")
	_ = GenKeyCmd.MarkFlagRequired("devid")
	_ = GenKeyCmd.MarkFlagRequired("prods")
}

// Copyright 2024 Google LLC
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

package instances

import (
	"fmt"
	"regexp"

	"internal/apiclient"
	"internal/client/hub"

	"github.com/spf13/cobra"
)

// CrtCmd to get a catalog items
var CrtCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an API Hub Instance",
	Long:  "Create an API Hub Instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		cmekPattern := `^projects/[^/]+/locations/[^/]+/keyRings/[^/]+/cryptoKeys/[^/]+$`
		if ok, err := regexp.MatchString(cmekPattern, cmekKeyName); err != nil || !ok {
			return fmt.Errorf("invalid pattern for cmekKeyName: %v", err)
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = hub.CreateInstance(id, cmekKeyName)
		return
	},
}

var cmekKeyName string

func init() {
	CrtCmd.Flags().StringVarP(&id, "id", "i",
		"", "API Hub Intance Id")
	CrtCmd.Flags().StringVarP(&cmekKeyName, "cmek-key-name", "k",
		"", "Cmek Key Name")

	CrtCmd.MarkFlagRequired("id")
	CrtCmd.MarkFlagRequired("cmek-key-name")
}

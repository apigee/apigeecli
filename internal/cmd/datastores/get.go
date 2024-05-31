// Copyright 2023 Google LLC
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

package datastores

import (
	"fmt"

	"internal/apiclient"
	"internal/client/datastores"

	"github.com/spf13/cobra"
)

// GetCmd to get a datastore
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a datastore connection",
	Long:  "Gets a datastore connection",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if name == "" && id == "" {
			return fmt.Errorf("id or name must be passed to the command")
		}
		if name != "" && id != "" {
			return fmt.Errorf("id and name cannot be passed together")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if id != "" {
			_, err = datastores.Get(id)
		} else {
			var respBody []byte
			if respBody, err = datastores.GetName(name); err != nil {
				return err
			}
			apiclient.PrettyPrint("application/json", respBody)
		}
		return
	},
}

func init() {
	GetCmd.Flags().StringVarP(&id, "id", "i",
		"", "Datastore UUID")
	GetCmd.Flags().StringVarP(&name, "name", "n",
		"", "Datastore display name")
}

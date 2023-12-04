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

package apidocs

import (
	"fmt"

	"internal/apiclient"
	"internal/client/apidocs"

	"github.com/spf13/cobra"
)

// ListCmd to list catalog items
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns the catalog items associated with a portal",
	Long:  "Returns the catalog items associated with a portal",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if siteid == "" {
			return fmt.Errorf("siteid is a mandatory parameter")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apidocs.List(siteid)
		return
	},
}

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

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/apps"
)

//GetCmd to get app
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get App in an Organization by App ID",
	Long:  "Returns the app profile for the specified app ID.",
	Args: func(cmd *cobra.Command, args []string) error {
		apiclient.SetApigeeOrg(org)
		if appID == "" && name == "" {
			return fmt.Errorf("pass either name or appId")
		}
		if appID != "" && name != "" {
			return fmt.Errorf("name and appId cannot be used together")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if appID != "" {
			_, err = apps.Get(appID)
			return
		}
		outBytes, err := apps.SearchApp(name)
		if err != nil {
			return err
		}
		//print the item
		return apiclient.PrettyPrint(outBytes)
	},
}

var appID string

func init() {

	GetCmd.Flags().StringVarP(&appID, "appId", "i",
		"", "Developer app id")
	GetCmd.Flags().StringVarP(&name, "name", "n",
		"", "Developer app name")
}

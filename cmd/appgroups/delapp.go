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

package appgroups

import (
	"internal/apiclient"
	"internal/client/appgroups"

	"github.com/spf13/cobra"
)

// DelAppCmd to get app
var DelAppCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an App in an AppGroup",
	Long:  "Delete an in an AppGroup",
	Args: func(cmd *cobra.Command, args []string) error {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = appgroups.DeleteApp(name, appName)
		return err
	},
}

func init() {
	DelAppCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the AppGroup")
	DelAppCmd.Flags().StringVarP(&appName, "app-name", "",
		"", "Name of the app")

	_ = DelAppCmd.MarkFlagRequired("name")
	_ = DelAppCmd.MarkFlagRequired("app-name")
}

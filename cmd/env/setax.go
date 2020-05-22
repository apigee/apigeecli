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

package env

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	environments "github.com/srinandan/apigeecli/client/env"
)

//Cmd to manage tracing of apis
var SetAxCmd = &cobra.Command{
	Use:   "setax",
	Short: "Set Analytics Agent role for a SA on an environment",
	Long:  "Set Analytics Agent role for a SA an Environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		apiclient.SetApigeeEnv(environment)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = environments.SetIAM(serviceAccountName, "analytics")
		if err != nil {
			return err
		}
		fmt.Printf("Service account %s granted access to Apigee Analytics Viewer role\n", serviceAccountName)
		return nil
	},
}

func init() {

	SetAxCmd.Flags().StringVarP(&serviceAccountName, "name", "n",
		"", "Service Account Name")

	_ = SetAxCmd.MarkFlagRequired("name")
}

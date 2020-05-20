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
	"regexp"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	environments "github.com/srinandan/apigeecli/client/env"
)

//SetCustCmd to manage custom roles for an env
var SetCustCmd = &cobra.Command{
	Use:   "setcustom",
	Short: "Set a custom role for a SA on an environment",
	Long:  "Set a custom role for a SA on an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		apiclient.SetApigeeEnv(environment)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		re := regexp.MustCompile(`projects\/([a-zA-Z0-9_-]+)\/roles\/([a-zA-Z0-9_-]+)`)
		result := re.FindString(role)
		if result == "" {
			return fmt.Errorf("custom role must be of the format projects/{project-id}/roles/{role-name}")
		}
		err = environments.SetIAM(serviceAccountName, role)
		if err != nil {
			return err
		}
		fmt.Printf("Service account %s, granted access to %s\n", serviceAccountName, role)
		return nil
	},
}

var role string

func init() {

	SetCustCmd.Flags().StringVarP(&serviceAccountName, "name", "n",
		"", "Service Account Name")
	SetCustCmd.Flags().StringVarP(&role, "role", "r",
		"", "Custom IAM role in the format projects/{project-id}/roles/{role}")

	_ = SetCustCmd.MarkFlagRequired("name")
	_ = SetCustCmd.MarkFlagRequired("role")
}

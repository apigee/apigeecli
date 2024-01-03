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

package preferences

import (
	"internal/apiclient"

	"github.com/spf13/cobra"
)

// SetCmd to get org details
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set default preferences for apigeecli",
	Long:  "Set default preferences for apigeecli",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = apiclient.WriteDefaultOrg(org); err != nil {
			return err
		}

		if err = apiclient.SetProxy(proxyURL); err != nil {
			return err
		}

		if err = apiclient.SetGithubURL(gitHubURL); err != nil {
			return err
		}

		if nocheck {
			if err = apiclient.SetNoCheck(nocheck); err != nil {
				return err
			}
		}

		if region != "" {
			if err = apiclient.WriteDefaultRegion(region); err != nil {
				return err
			}
		}

		return nil
	},
}

var (
	org, region, proxyURL, gitHubURL string
	usestage, nocheck                bool
)

func init() {
	SetCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	SetCmd.Flags().BoolVarP(&usestage, "staging", "s",
		false, "Use Apigee staging; format: -s=true")

	SetCmd.Flags().StringVarP(&proxyURL, "proxy", "p",
		"", "Use http proxy before contacting the control plane")

	SetCmd.Flags().BoolVarP(&nocheck, "nocheck", "",
		false, "Don't check for newer versions of cmd")

	SetCmd.Flags().StringVarP(&gitHubURL, "github", "g",
		"", "On premises Github URL")

	SetCmd.Flags().StringVarP(&region, "region", "r",
		"", "Apigee control plane region")
}

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

package securityprofiles

import (
	"fmt"
	"internal/apiclient"
	"internal/client/securityprofiles"
	"time"

	"github.com/spf13/cobra"
)

// ScoreCmd fetches the Advanced API Security risk assessment score for a proxy
var ScoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Get Advanced API Security risk assessment score for an API proxy",
	Long: "Get the Advanced API Security risk assessment score for a specific API proxy " +
		"in an environment. Useful as a CI/CD guardrail to block deployments when the score is too low.",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if name == "" {
			return fmt.Errorf("security profile name cannot be empty")
		}
		if proxyName == "" {
			return fmt.Errorf("proxy name cannot be empty")
		}
		if scoreStartTime != "" {
			if _, err = time.Parse(time.RFC3339, scoreStartTime); err != nil {
				return fmt.Errorf("invalid format for start-time: %v", err)
			}
		}
		if scoreEndTime != "" {
			if _, err = time.Parse(time.RFC3339, scoreEndTime); err != nil {
				return fmt.Errorf("invalid format for end-time: %v", err)
			}
		}
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = securityprofiles.GetProxyScore(name, proxyName, scoreStartTime, scoreEndTime)
		return
	},
}

var (
	proxyName        string
	scoreStartTime   string
	scoreEndTime     string
)

func init() {
	ScoreCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the security profile")
	ScoreCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")
	ScoreCmd.Flags().StringVarP(&proxyName, "proxy", "p",
		"", "API proxy name to retrieve the risk assessment score for")
	ScoreCmd.Flags().StringVarP(&scoreStartTime, "start-time", "",
		"", "Inclusive start of the interval in RFC3339 format; default is 24 hours ago")
	ScoreCmd.Flags().StringVarP(&scoreEndTime, "end-time", "",
		"", "Exclusive end of the interval in RFC3339 format; default is current timestamp")

	_ = ScoreCmd.MarkFlagRequired("name")
	_ = ScoreCmd.MarkFlagRequired("env")
	_ = ScoreCmd.MarkFlagRequired("proxy")
}

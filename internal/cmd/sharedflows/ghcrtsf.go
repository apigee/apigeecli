// Copyright 2021 Google LLC
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

package sharedflows

import (
	"fmt"
	"internal/apiclient"
	"internal/client/sharedflows"
	"internal/clilog"
	"os"
	"regexp"

	proxybundle "internal/bundlegen/proxybundle"

	"github.com/spf13/cobra"
)

// GhCreateCmd create an api from a github repo
var GhCreateCmd = &cobra.Command{
	Use:     "github",
	Aliases: []string{"gh"},
	Short:   "Creates a sharedflow from a GitHub repo",
	Long:    "Creates a sharedflow from a GitHub repo. Check apigeecli prefs for GH on-prem options",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		re := regexp.MustCompile(`(\w+)?\/sharedflowbundle$`)
		if ok := re.Match([]byte(ghPath)); !ok {
			return fmt.Errorf("github path must end with /sharedflowbundle")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if os.Getenv("GITHUB_TOKEN") == "" {
			clilog.Debug.Println("github token is not set as an env var. Running unauthenticated")
		}
		if err = proxybundle.GitHubImportBundle(ghOwner, ghRepo, ghPath, true); err != nil {
			proxybundle.SharedflowCleanUp()
			return err
		}

		_, err = sharedflows.Create(name, bundleName, space)
		proxybundle.SharedflowCleanUp()
		return err
	},
}

const bundleName = "sharedflowbundle.zip"

var ghOwner, ghRepo, ghPath string

func init() {
	GhCreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Sharedflow name")
	GhCreateCmd.Flags().StringVarP(&ghOwner, "owner", "u",
		"", "The github organization or username. ex: In https://github.com/apigee, apigee is the owner name")
	GhCreateCmd.Flags().StringVarP(&space, "space", "",
		"", "Apigee Space to associate to")
	GhCreateCmd.Flags().StringVarP(&ghRepo, "repo", "",
		"", "The github repo name. ex: https://github.com/apigee/api-platform-samples, api-platform-samples is the repo")
	GhCreateCmd.Flags().StringVarP(&ghPath, "sf-path", "p",
		"", "The path in the repo to the sharedflowbundle folder. ex: sample-proxies/security/sharedflowbundle")

	_ = GhCreateCmd.MarkFlagRequired("name")
	_ = GhCreateCmd.MarkFlagRequired("owner")
	_ = GhCreateCmd.MarkFlagRequired("repo")
	_ = GhCreateCmd.MarkFlagRequired("sf-path")
}

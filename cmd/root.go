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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/cmd/apis"
	"github.com/srinandan/apigeecli/cmd/apps"
	cache "github.com/srinandan/apigeecli/cmd/cache"
	"github.com/srinandan/apigeecli/cmd/datacollectors"
	"github.com/srinandan/apigeecli/cmd/developers"
	"github.com/srinandan/apigeecli/cmd/env"
	"github.com/srinandan/apigeecli/cmd/envgroup"
	"github.com/srinandan/apigeecli/cmd/envoy"
	"github.com/srinandan/apigeecli/cmd/eptattachment"
	flowhooks "github.com/srinandan/apigeecli/cmd/flowhooks"
	"github.com/srinandan/apigeecli/cmd/iam"
	"github.com/srinandan/apigeecli/cmd/instances"
	"github.com/srinandan/apigeecli/cmd/keyaliases"
	"github.com/srinandan/apigeecli/cmd/keystores"
	"github.com/srinandan/apigeecli/cmd/kvm"
	"github.com/srinandan/apigeecli/cmd/ops"
	"github.com/srinandan/apigeecli/cmd/org"
	"github.com/srinandan/apigeecli/cmd/overrides"
	"github.com/srinandan/apigeecli/cmd/preferences"
	"github.com/srinandan/apigeecli/cmd/products"
	"github.com/srinandan/apigeecli/cmd/projects"
	"github.com/srinandan/apigeecli/cmd/references"
	res "github.com/srinandan/apigeecli/cmd/res"
	"github.com/srinandan/apigeecli/cmd/sharedflows"
	"github.com/srinandan/apigeecli/cmd/sync"
	targetservers "github.com/srinandan/apigeecli/cmd/targetservers"
	"github.com/srinandan/apigeecli/cmd/token"
)

//RootCmd to manage apigeecli
var RootCmd = &cobra.Command{
	Use:   "apigeecli",
	Short: "Utility to work with Apigee APIs.",
	Long:  "This command lets you interact with Apigee APIs.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		apiclient.SetServiceAccount(serviceAccount)
		apiclient.SetApigeeToken(accessToken)

		if !disableCheck {
			if ok, _ := apiclient.TestAndUpdateLastCheck(); !ok {
				latestVersion, _ := getLatestVersion()
				if cmd.Version != latestVersion {
					fmt.Printf("You are using %s, the latest version %s is available for download\n", cmd.Version, latestVersion)
				}
			}
		}

		_ = apiclient.SetAccessToken()

		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var accessToken, serviceAccount string
var disableCheck bool

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&accessToken, "token", "t",
		"", "Google OAuth Token")

	RootCmd.PersistentFlags().StringVarP(&serviceAccount, "account", "a",
		"", "Path Service Account private key in JSON")

	RootCmd.PersistentFlags().BoolVarP(&disableCheck, "disable-check", "",
		false, "Disable check for newer versions")

	RootCmd.AddCommand(apis.Cmd)
	RootCmd.AddCommand(org.Cmd)
	RootCmd.AddCommand(sync.Cmd)
	RootCmd.AddCommand(envgroup.Cmd)
	RootCmd.AddCommand(env.Cmd)
	RootCmd.AddCommand(products.Cmd)
	RootCmd.AddCommand(datacollectors.Cmd)
	RootCmd.AddCommand(developers.Cmd)
	RootCmd.AddCommand(apps.Cmd)
	RootCmd.AddCommand(sharedflows.Cmd)
	RootCmd.AddCommand(kvm.Cmd)
	RootCmd.AddCommand(flowhooks.Cmd)
	RootCmd.AddCommand(targetservers.Cmd)
	RootCmd.AddCommand(token.Cmd)
	RootCmd.AddCommand(keystores.Cmd)
	RootCmd.AddCommand(keyaliases.Cmd)
	RootCmd.AddCommand(cache.Cmd)
	RootCmd.AddCommand(references.Cmd)
	RootCmd.AddCommand(res.Cmd)
	RootCmd.AddCommand(projects.Cmd)
	RootCmd.AddCommand(iam.Cmd)
	RootCmd.AddCommand(envoy.Cmd)
	RootCmd.AddCommand(instances.Cmd)
	RootCmd.AddCommand(ops.Cmd)
	RootCmd.AddCommand(preferences.Cmd)
	RootCmd.AddCommand(overrides.Cmd)
	RootCmd.AddCommand(eptattachment.Cmd)
}

func initConfig() {
	var skipLogInfo = true
	var skipCache bool

	if os.Getenv("APIGEECLI_SKIPLOG") == "false" {
		skipLogInfo = false
	}

	skipCache, _ = strconv.ParseBool(os.Getenv("APIGEECLI_SKIPCACHE"))

	apiclient.NewApigeeClient(apiclient.ApigeeClientOptions{
		SkipCheck:   true,
		PrintOutput: true,
		SkipLogInfo: skipLogInfo,
		SkipCache:   skipCache,
	})
}

// GetRootCmd returns the root of the cobra command-tree.
func GetRootCmd() *cobra.Command {
	return RootCmd
}

func getLatestVersion() (version string, err error) {
	var req *http.Request
	const endpoint = "https://api.github.com/repos/srinandan/apigeecli/releases/latest"

	client := &http.Client{}
	contentType := "application/json"

	req, err = http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", result["tag_name"]), nil
}

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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"internal/apiclient"

	"internal/clilog"

	"github.com/apigee/apigeecli/cmd/apicategories"
	"github.com/apigee/apigeecli/cmd/apidocs"
	"github.com/apigee/apigeecli/cmd/apis"
	"github.com/apigee/apigeecli/cmd/appgroups"
	"github.com/apigee/apigeecli/cmd/apps"
	cache "github.com/apigee/apigeecli/cmd/cache"
	"github.com/apigee/apigeecli/cmd/datacollectors"
	"github.com/apigee/apigeecli/cmd/datastores"
	"github.com/apigee/apigeecli/cmd/developers"
	"github.com/apigee/apigeecli/cmd/env"
	"github.com/apigee/apigeecli/cmd/envgroup"
	"github.com/apigee/apigeecli/cmd/eptattachment"
	flowhooks "github.com/apigee/apigeecli/cmd/flowhooks"
	"github.com/apigee/apigeecli/cmd/iam"
	"github.com/apigee/apigeecli/cmd/instances"
	"github.com/apigee/apigeecli/cmd/keyaliases"
	"github.com/apigee/apigeecli/cmd/keystores"
	"github.com/apigee/apigeecli/cmd/kvm"
	"github.com/apigee/apigeecli/cmd/ops"
	"github.com/apigee/apigeecli/cmd/org"
	"github.com/apigee/apigeecli/cmd/preferences"
	"github.com/apigee/apigeecli/cmd/products"
	"github.com/apigee/apigeecli/cmd/projects"
	"github.com/apigee/apigeecli/cmd/references"
	res "github.com/apigee/apigeecli/cmd/res"
	"github.com/apigee/apigeecli/cmd/securityprofiles"
	"github.com/apigee/apigeecli/cmd/sharedflows"
	"github.com/apigee/apigeecli/cmd/sync"
	targetservers "github.com/apigee/apigeecli/cmd/targetservers"
	"github.com/apigee/apigeecli/cmd/token"
	"github.com/spf13/cobra"
)

// RootCmd to manage apigeecli
var RootCmd = &cobra.Command{
	Use:   "apigeecli",
	Short: "Utility to work with Apigee APIs.",
	Long:  "This command lets you interact with Apigee APIs.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if metadataToken && defaultToken {
			return fmt.Errorf("metadata-token and default-token cannot be used together")
		}
		if defaultToken && (serviceAccount != "" || accessToken != "") {
			return fmt.Errorf("default-token cannot be used with token or account flags")
		}
		if metadataToken && (serviceAccount != "" || accessToken != "") {
			return fmt.Errorf("metadata-token cannot be used with token or account flags")
		}

		if serviceAccount != "" && accessToken != "" {
			return fmt.Errorf("token and account flags cannot be used together")
		}

		if !disableCheck {
			if ok, _ := apiclient.TestAndUpdateLastCheck(); !ok {
				latestVersion, _ := getLatestVersion()
				if cmd.Version == "" {
					clilog.Debug.Println("apigeecli wasn't built with a valid Version tag.")
				} else if latestVersion != "" && cmd.Version != latestVersion {
					clilog.Info.Printf("You are using %s, the latest version %s "+
						"is available for download\n", cmd.Version, latestVersion)
				}
			}
		}

		if !metadataToken && !defaultToken {
			apiclient.SetServiceAccount(serviceAccount)
			apiclient.SetApigeeToken(accessToken)
		}

		if metadataToken {
			return apiclient.GetMetadataAccessToken()
		}

		if defaultToken {
			return apiclient.GetDefaultAccessToken()
		}

		_ = apiclient.SetAccessToken()

		return nil
	},
	SilenceUsage:  getUsageFlag(),
	SilenceErrors: getErrorsFlag(),
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		clilog.Error.Println(err)
	}
}

var (
	accessToken, serviceAccount                                      string
	disableCheck, printOutput, noOutput, metadataToken, defaultToken bool
)

const ENABLED = "true"

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&accessToken, "token", "t",
		"", "Google OAuth Token")

	RootCmd.PersistentFlags().StringVarP(&serviceAccount, "account", "a",
		"", "Path Service Account private key in JSON")

	RootCmd.PersistentFlags().BoolVarP(&disableCheck, "disable-check", "",
		false, "Disable check for newer versions")

	RootCmd.PersistentFlags().BoolVarP(&printOutput, "print-output", "",
		true, "Control printing of info log statements")

	RootCmd.PersistentFlags().BoolVarP(&noOutput, "no-output", "",
		false, "Disable printing all statements to stdout")

	RootCmd.PersistentFlags().BoolVarP(&metadataToken, "metadata-token", "",
		false, "Metadata OAuth2 access token")

	RootCmd.PersistentFlags().BoolVarP(&defaultToken, "default-token", "",
		false, "Use Google default application credentials access token")

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
	RootCmd.AddCommand(instances.Cmd)
	RootCmd.AddCommand(ops.Cmd)
	RootCmd.AddCommand(preferences.Cmd)
	RootCmd.AddCommand(eptattachment.Cmd)
	RootCmd.AddCommand(appgroups.Cmd)
	RootCmd.AddCommand(apidocs.Cmd)
	RootCmd.AddCommand(apicategories.Cmd)
	RootCmd.AddCommand(datastores.Cmd)
	RootCmd.AddCommand(securityprofiles.Cmd)
}

func initConfig() {
	debug := false
	var skipCache bool

	if os.Getenv("APIGEECLI_DEBUG") == ENABLED {
		debug = true
	}

	skipCache, _ = strconv.ParseBool(os.Getenv("APIGEECLI_SKIPCACHE"))

	if noOutput {
		printOutput = noOutput
	}

	apiclient.NewApigeeClient(apiclient.ApigeeClientOptions{
		TokenCheck:  true,
		PrintOutput: printOutput,
		NoOutput:    noOutput,
		DebugLog:    debug,
		SkipCache:   skipCache,
	})

	if os.Getenv("APIGEECLI_ENABLE_RATELIMIT") == ENABLED {
		clilog.Debug.Println("APIGEECLI_RATELIMIT is enabled")
		apiclient.SetRate(apiclient.ApigeeAPI)
	}
}

// GetRootCmd returns the root of the cobra command-tree.
func GetRootCmd() *cobra.Command {
	return RootCmd
}

func getLatestVersion() (version string, err error) {
	var req *http.Request
	const endpoint = "https://api.github.com/repos/apigee/apigeecli/releases/latest"

	client := &http.Client{}
	contentType := "application/json"

	ctx := context.Background()
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", err
	}

	if result["tag_name"] == "" {
		clilog.Debug.Println("Unable to determine latest tag, skipping this information")
		return "", nil
	}
	return fmt.Sprintf("%s", result["tag_name"]), nil
}

// getUsageFlag
func getUsageFlag() bool {
	return os.Getenv("APIGEECLI_NO_USAGE") == ENABLED
}

// getErrorsFlag
func getErrorsFlag() bool {
	return os.Getenv("APIGEECLI_NO_ERRORS") == ENABLED
}

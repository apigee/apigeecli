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

	"internal/cmd/apicategories"
	"internal/cmd/apidocs"
	"internal/cmd/apihub"
	"internal/cmd/apis"
	"internal/cmd/appgroups"
	"internal/cmd/apps"
	cache "internal/cmd/cache"
	"internal/cmd/datacollectors"
	"internal/cmd/datastores"
	"internal/cmd/developers"
	"internal/cmd/env"
	"internal/cmd/envgroup"
	"internal/cmd/eptattachment"
	flowhooks "internal/cmd/flowhooks"
	"internal/cmd/iam"
	"internal/cmd/instances"
	"internal/cmd/keyaliases"
	"internal/cmd/keystores"
	"internal/cmd/kvm"
	"internal/cmd/observe"
	"internal/cmd/ops"
	"internal/cmd/org"
	"internal/cmd/preferences"
	"internal/cmd/products"
	"internal/cmd/projects"
	"internal/cmd/references"
	res "internal/cmd/res"
	"internal/cmd/securityprofiles"
	"internal/cmd/sharedflows"
	"internal/cmd/sites"
	"internal/cmd/sync"
	targetservers "internal/cmd/targetservers"
	"internal/cmd/token"

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

		apiclient.SetAPI(api)

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
	accessToken, serviceAccount                                                  string
	disableCheck, printOutput, noOutput, metadataToken, defaultToken, noWarnings bool
	api                                                                          apiclient.API
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

	RootCmd.PersistentFlags().BoolVarP(&noWarnings, "no-warnings", "",
		false, "Disable printing warnings to stderr")

	RootCmd.PersistentFlags().BoolVarP(&metadataToken, "metadata-token", "",
		false, "Metadata OAuth2 access token")

	RootCmd.PersistentFlags().BoolVarP(&defaultToken, "default-token", "",
		false, "Use Google default application credentials access token")

	RootCmd.PersistentFlags().Var(&api, "api", "Sets the control plane API. Must be one of prod, autopush "+
		"or staging; default is prod")

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
	RootCmd.AddCommand(sites.Cmd)
	RootCmd.AddCommand(apihub.Cmd)
	RootCmd.AddCommand(observe.Cmd)
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
		NoWarnings:  noWarnings,
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

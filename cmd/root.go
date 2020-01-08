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
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/clilog"
	"github.com/srinandan/apigeecli/cmd/apis"
	"github.com/srinandan/apigeecli/cmd/apps"
	cache "github.com/srinandan/apigeecli/cmd/cache"
	"github.com/srinandan/apigeecli/cmd/developers"
	"github.com/srinandan/apigeecli/cmd/env"
	flowhooks "github.com/srinandan/apigeecli/cmd/flowhooks"
	"github.com/srinandan/apigeecli/cmd/iam"
	"github.com/srinandan/apigeecli/cmd/keyaliases"
	"github.com/srinandan/apigeecli/cmd/keystores"
	"github.com/srinandan/apigeecli/cmd/kvm"
	"github.com/srinandan/apigeecli/cmd/org"
	"github.com/srinandan/apigeecli/cmd/products"
	"github.com/srinandan/apigeecli/cmd/projects"
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
		clilog.Init(apiclient.IsSkipLogInfo())
		return apiclient.SetAccessToken()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(apiclient.SkipLogInfo(), "log", "l",
		false, "Log Information")

	RootCmd.PersistentFlags().StringVarP(apiclient.GetApigeeTokenP(), "token", "t",
		"", "Google OAuth Token")

	RootCmd.PersistentFlags().StringVarP(apiclient.GetServiceAccountP(), "account", "a",
		"", "Path Service Account private key in JSON")

	RootCmd.PersistentFlags().BoolVar(apiclient.SkipCache(), "skipCache",
		false, "Skip caching Google OAuth Token")

	RootCmd.PersistentFlags().BoolVar(apiclient.SkipCheck(), "skipCheck",
		true, "Skip checking expiry for Google OAuth Token")

	RootCmd.AddCommand(apis.Cmd)
	RootCmd.AddCommand(org.Cmd)
	RootCmd.AddCommand(sync.Cmd)
	RootCmd.AddCommand(env.Cmd)
	RootCmd.AddCommand(products.Cmd)
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
	RootCmd.AddCommand(res.Cmd)
	RootCmd.AddCommand(projects.Cmd)
	RootCmd.AddCommand(iam.Cmd)
}

func initConfig() {

}

// GetRootCmd returns the root of the cobra command-tree.
func GetRootCmd() *cobra.Command {
	return RootCmd
}

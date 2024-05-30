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
	"internal/apiclient"

	environments "internal/client/env"

	"github.com/spf13/cobra"
)

// CrtTraceOverridesCmd to manage tracing of apis
var CrtTraceOverridesCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Distributed Trace config override",
	Long:  "Create a new Distributed Trace config override",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = environments.CreateTraceOverrides(apiProxy, exporter, endpoint, sampler, sampleRate)
		return
	},
}

var apiProxy string

func init() {
	CrtTraceOverridesCmd.Flags().StringVarP(&apiProxy, "name", "n",
		"", "API Proxy name")
	CrtTraceOverridesCmd.Flags().StringVarP(&exporter, "exporter", "x",
		"", "Trace exporter can be JAEGER or CLOUD_TRACE")
	CrtTraceOverridesCmd.Flags().StringVarP(&endpoint, "endpoint-uri", "p",
		"", "Trace endpoint, used only with JAEGER")
	CrtTraceOverridesCmd.Flags().StringVarP(&sampler, "sampler", "s",
		"PROBABILITY", "Sampler can be set to PROBABILITY or OFF")
	CrtTraceOverridesCmd.Flags().StringVarP(&sampleRate, "rate", "",
		"", "Sampler Rate")

	_ = CrtTraceOverridesCmd.MarkFlagRequired("exporter")
	_ = CrtTraceOverridesCmd.MarkFlagRequired("name")
}

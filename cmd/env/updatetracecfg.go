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
	"github.com/apigee/apigeecli/apiclient"
	environments "github.com/apigee/apigeecli/client/env"
	"github.com/spf13/cobra"
)

//UpdateTraceConfigCmd to manage tracing of apis
var UpdateTraceConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Distributed Trace config for the environment",
	Long:  "Update Distributed Trace config for the environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = environments.UpdateTraceConfig(exporter, endpoint, sampler, sample_rate)
		return
	},
}

var exporter, endpoint, sampler, sample_rate string

func init() {

	UpdateTraceConfigCmd.Flags().StringVarP(&exporter, "exporter", "x",
		"", "Trace exporter can be JAEGER or CLOUD_TRACE")
	UpdateTraceConfigCmd.Flags().StringVarP(&endpoint, "endpoint", "p",
		"", "Trace endpoint, used only with JAEGER")
	UpdateTraceConfigCmd.Flags().StringVarP(&sampler, "sampler", "s",
		"PROBABILITY", "Sampler can be set to PROBABILITY or OFF")
	UpdateTraceConfigCmd.Flags().StringVarP(&sample_rate, "rate", "r",
		"", "Sampler Rate")

	_ = UpdateTraceConfigCmd.MarkFlagRequired("exporter")
}

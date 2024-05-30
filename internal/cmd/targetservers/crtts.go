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

package targetservers

import (
	"fmt"
	"strconv"

	"internal/apiclient"

	"internal/client/targetservers"

	"github.com/spf13/cobra"
)

// CreateCmd to create target servers
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Target Server",
	Long:  "Create a Target Server",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		if tlsenabled != "" {
			if _, err := strconv.ParseBool(tlsenabled); err != nil {
				return fmt.Errorf("tlsenabled must be set to  true or false")
			}
		}
		if tlsenforce != "" {
			if _, err := strconv.ParseBool(tlsenforce); err != nil {
				return fmt.Errorf("tlsenforce must be set to  true or false")
			}
		}
		if clientAuthEnabled != "" {
			if _, err := strconv.ParseBool(clientAuthEnabled); err != nil {
				return fmt.Errorf("clientAuthEnabled must be set to  true or false")
			}
		}
		if ignoreValidationErrors != "" {
			if _, err := strconv.ParseBool(ignoreValidationErrors); err != nil {
				return fmt.Errorf("ignoreValidationErrors must be set to  true or false")
			}
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = targetservers.Create(name,
			description,
			host,
			port,
			enable,
			protocol,
			keyStore, keyAlias, trustStore,
			tlsenabled, tlsenforce, clientAuthEnabled,
			ignoreValidationErrors)
		return err
	},
}

var (
	tlsenabled, tlsenforce, clientAuthEnabled, description, host, keyStore, keyAlias string
	trustStore, protocol, ignoreValidationErrors                                     string
	enable                                                                           bool
	port                                                                             int
)

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the targetserver")
	CreateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the Target Server")
	CreateCmd.Flags().StringVarP(&host, "host", "s",
		"", "Host name of the target")
	CreateCmd.Flags().BoolVarP(&enable, "enable", "b",
		true, "Enabling/disabling a TargetServer")

	CreateCmd.Flags().StringVarP(&keyStore, "keystore", "",
		"", "Key store for the target server")
	CreateCmd.Flags().StringVarP(&keyAlias, "keyalias", "",
		"", "Key alias for the target server")
	CreateCmd.Flags().StringVarP(&trustStore, "truststore", "",
		"", "Trust store for the target server")

	CreateCmd.Flags().StringVarP(&tlsenabled, "tls", "",
		"", "Enable TLS for the target server")
	CreateCmd.Flags().StringVarP(&tlsenforce, "tlsenforce", "",
		"", "Enforce TLS for the target server")
	CreateCmd.Flags().StringVarP(&clientAuthEnabled, "client-auth", "c",
		"", "Enable mTLS for the target server")
	CreateCmd.Flags().StringVarP(&ignoreValidationErrors, "ignore-err", "i",
		"", "Ignore TLS validation errors for the target server")

	CreateCmd.Flags().IntVarP(&port, "port", "p",
		-1, "port number")
	CreateCmd.Flags().StringVarP(&protocol, "protocol", "",
		"HTTP", "Protocol for a TargetServer; default is HTTP")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("host")
}

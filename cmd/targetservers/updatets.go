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

// UpdateCmd to get target servers
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a Target Server",
	Long:  "Update a Target Server",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		if tlsenabled != "" {
			if _, err := strconv.ParseBool(tlsenabled); err != nil {
				return fmt.Errorf("tlsenabled must be set to  true or false")
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
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		_, err = targetservers.Update(name,
			description,
			host,
			port,
			enable,
			protocol,
			keyStore, keyAlias, trustStore,
			tlsenabled, clientAuthEnabled,
			ignoreValidationErrors)
		return err
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the targetserver")
	UpdateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the Target Server")
	UpdateCmd.Flags().StringVarP(&host, "host", "s",
		"", "Host name of the target")
	UpdateCmd.Flags().BoolVarP(&enable, "enable", "b",
		true, "Enabling/disabling a TargetServer")

	UpdateCmd.Flags().StringVarP(&keyStore, "keyStore", "",
		"", "Key store for the target server")
	UpdateCmd.Flags().StringVarP(&keyAlias, "keyAlias", "",
		"", "Key alias for the target server")
	UpdateCmd.Flags().StringVarP(&trustStore, "trustStore", "",
		"", "Trust store for the target server")

	UpdateCmd.Flags().StringVarP(&tlsenabled, "tls", "",
		"", "Enable TLS for the target server")
	UpdateCmd.Flags().StringVarP(&clientAuthEnabled, "clientAuth", "c",
		"", "Enable mTLS for the target server")
	UpdateCmd.Flags().StringVarP(&ignoreValidationErrors, "ignoreErr", "i",
		"", "Ignore TLS validation errors for the target server")

	UpdateCmd.Flags().IntVarP(&port, "port", "p",
		-1, "port number")
	UpdateCmd.Flags().StringVarP(&protocol, "protocol", "",
		"", "Protocol for a TargetServer")

	_ = UpdateCmd.MarkFlagRequired("name")
}

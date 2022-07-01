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

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/targetservers"
	"github.com/spf13/cobra"
)

//UpdateCmd to get target servers
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a Target Server",
	Long:  "Update a Target Server",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		if sslinfo != "" {
			if _, err = strconv.ParseBool(sslinfo); err != nil {
				return fmt.Errorf("Invalid value for sslinfo. Must be set to true or false")
			}
		}
		if enable != "" {
			if _, err = strconv.ParseBool(enable); err != nil {
				return fmt.Errorf("Invalid value for enable. Must be set to true or false")
			}
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = targetservers.Update(name, description, host, port, enable, grpc, keyStore, keyAlias, sslinfo, tlsenabled, clientAuthEnabled, ignoreValidationErrors)
		return
	},
}

func init() {

	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the targetserver")
	UpdateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the Target Server")
	UpdateCmd.Flags().StringVarP(&host, "host", "s",
		"", "Host name of the target")
	UpdateCmd.Flags().StringVarP(&enable, "enable", "b",
		"", "Enabling/disabling a TargetServer")
	UpdateCmd.Flags().BoolVarP(&grpc, "grpc", "g",
		false, "Enable target server for gRPC")

	UpdateCmd.Flags().StringVarP(&keyStore, "keyStore", "",
		"", "Key store for the target server; must be used with sslinfo")
	UpdateCmd.Flags().StringVarP(&keyAlias, "keyAlias", "",
		"", "Key alias for the target server; must be used with sslinfo")
	UpdateCmd.Flags().StringVarP(&sslinfo, "sslinfo", "",
		"", "Enable SSL Info on the target server")
	UpdateCmd.Flags().BoolVarP(&tlsenabled, "tls", "",
		false, "Enable TLS for the target server; must be used with sslinfo")
	UpdateCmd.Flags().BoolVarP(&clientAuthEnabled, "clientAuth", "c",
		false, "Enable mTLS for the target server; must be used with sslinfo")
	UpdateCmd.Flags().BoolVarP(&ignoreValidationErrors, "ignoreErr", "i",
		false, "Ignore TLS validation errors for the target server; must be used with sslinfo")

	UpdateCmd.Flags().IntVarP(&port, "port", "p",
		-1, "port number")

	_ = UpdateCmd.MarkFlagRequired("name")

}

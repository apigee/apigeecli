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
	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/targetservers"
	"github.com/spf13/cobra"
)

//CreateCmd to create target servers
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Target Server",
	Long:  "Create a Target Server",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = targetservers.Create(name, description, host, port, enable, grpc, keyStore, keyAlias, sslinfo, tlsenabled, clientAuthEnabled, ignoreValidationErrors)
		return

	},
}

var description, host, keyStore, keyAlias string
var enable, grpc, sslinfo, tlsenabled, clientAuthEnabled, ignoreValidationErrors bool
var port int

func init() {

	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the targetserver")
	CreateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the Target Server")
	CreateCmd.Flags().StringVarP(&host, "host", "s",
		"", "Host name of the target")
	CreateCmd.Flags().BoolVarP(&enable, "enable", "b",
		true, "Enabling/disabling a TargetServer")
	CreateCmd.Flags().BoolVarP(&grpc, "grpc", "g",
		false, "Enable target server for gRPC")

	CreateCmd.Flags().StringVarP(&keyStore, "keyStore", "",
		"", "Key store for the target server; must be used with sslinfo")
	CreateCmd.Flags().StringVarP(&keyAlias, "keyAlias", "",
		"", "Key alias for the target server; must be used with sslinfo")
	CreateCmd.Flags().BoolVarP(&sslinfo, "sslinfo", "",
		false, "Enable SSL Info on the target server")
	CreateCmd.Flags().BoolVarP(&tlsenabled, "tls", "",
		false, "Enable TLS for the target server; must be used with sslinfo")
	CreateCmd.Flags().BoolVarP(&clientAuthEnabled, "clientAuth", "c",
		false, "Enable mTLS for the target server; must be used with sslinfo")
	CreateCmd.Flags().BoolVarP(&ignoreValidationErrors, "ignoreErr", "i",
		false, "Ignore TLS validation errors for the target server; must be used with sslinfo")

	CreateCmd.Flags().IntVarP(&port, "port", "p",
		-1, "port number")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("host")
}

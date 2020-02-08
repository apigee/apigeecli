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
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/targetservers"
)

//CreateCmd to create target servers
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Target Server",
	Long:  "Create a Target Server",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		apiclient.SetApigeeEnv(env)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = targetservers.Create(name, description, host, port, enable)
		return

	},
}

var description, host string
var enable bool
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
	CreateCmd.Flags().IntVarP(&port, "port", "p",
		80, "port number")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("host")
}

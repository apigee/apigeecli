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

package instances

import (
	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/instances"
	"github.com/spf13/cobra"
)

//ReserveNatCmd activates NAT for an Apigee instance
var ReserveNatCmd = &cobra.Command{
	Use:   "reserve",
	Short: "Reserve a NAT IP for an Apigeee instance",
	Long:  "Reserve a NAT IP for an Apigeee instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = instances.ReserveNatIP(name, natid)
		return
	},
}

func init() {
	ReserveNatCmd.Flags().StringVarP(&natid, "natid", "i",
		"", "NAT identifier")

	_ = ReserveNatCmd.MarkFlagRequired("natid")
}

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
	"github.com/spf13/cobra"
)

//NatCmd to manage instances
var NatCmd = &cobra.Command{
	Use:   "nat",
	Short: "Manage NAT IPs for Apigee instances",
	Long:  "Manage NAT IPs for Apigee instances",
}

var natid string

func init() {

	NatCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	NatCmd.PersistentFlags().StringVarP(&name, "name", "n",
		"", "Apigee instance name")

	NatCmd.AddCommand(ReserveNatCmd)
	NatCmd.AddCommand(ActivateNatCmd)
	NatCmd.AddCommand(ListNatCmd)
	NatCmd.AddCommand(DeleteNatCmd)

	_ = NatCmd.MarkPersistentFlagRequired("name")
}

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

package developers

import (
	"internal/apiclient"

	"internal/client/developers"

	"github.com/spf13/cobra"
)

// CreditCmd to list developer subscriptions
var CreditCmd = &cobra.Command{
	Use:   "credit",
	Short: "Credits the account balance for the developer",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	Long: "Credits the account balance for the developer",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = developers.Credit(email, credit)
		return
	},
}

var credit string

func init() {
	CreditCmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")
	CreditCmd.Flags().StringVarP(&credit, "credit", "",
		"", "Credit transaction payload in json")

	_ = CreditCmd.MarkFlagRequired("email")
	_ = CreditCmd.MarkFlagRequired("credit")
}

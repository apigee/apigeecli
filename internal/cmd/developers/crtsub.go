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

// CreateSubCmd to create developer subscription
var CreateSubCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a subcription for a developer",
	Long:  "Create a developer",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = developers.CreateSubscription(email, name, apiproduct, startTime, endTime)
		return
	},
}

var apiproduct, startTime, endTime string

func init() {
	CreateSubCmd.Flags().StringVarP(&email, "email", "e",
		"", "The developer's email")
	CreateSubCmd.Flags().StringVarP(&name, "name", "n",
		"", "The subscription name")
	CreateSubCmd.Flags().StringVarP(&apiproduct, "apiproduct", "p",
		"", "The name of the API Product")
	CreateSubCmd.Flags().StringVarP(&startTime, "start", "s",
		"", "Subscription start time")
	CreateSubCmd.Flags().StringVarP(&endTime, "end", "d",
		"", "Subscription end time")

	_ = CreateSubCmd.MarkFlagRequired("email")
	_ = CreateSubCmd.MarkFlagRequired("apiproduct")
	_ = CreateSubCmd.MarkFlagRequired("start")
	_ = CreateSubCmd.MarkFlagRequired("end")
}

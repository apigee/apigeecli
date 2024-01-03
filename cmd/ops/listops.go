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

package ops

import (
	"fmt"

	"internal/apiclient"

	"internal/client/operations"

	"github.com/spf13/cobra"
)

// ListCmd to list envs
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List operations in an Apigee Org",
	Long:  "List operations in an Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if state == "" && completeState != "" {
			return fmt.Errorf("both state and completeState must be passed")
		}
		if completeState != "Success" && completeState != "Failed" && completeState != "Both" && completeState != "" {
			return fmt.Errorf("completeState must be oneOf: Success, Failed, Both or empty")
		}
		if state != "IN_PROGRESS" && state != "FINISHED" && state != "ERROR" {
			return fmt.Errorf("state must be oneOf IN_PROGRESS, FINISHED or ERROR")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if completeState == "Success" {
			_, err = operations.List(state, operations.Success)
		} else if completeState == "Failed" {
			_, err = operations.List(state, operations.Failed)
		} else {
			_, err = operations.List(state, operations.Both)
		}
		return err
	},
}

var (
	state         string
	completeState string
)

func init() {
	ListCmd.Flags().StringVarP(&state, "state", "s",
		"", "filter by operation state: FINISHED, ERROR, IN_PROGRESS")
	ListCmd.Flags().StringVarP(&completeState, "completeState", "c",
		"", "filter by operation compeleted state: Success, Failed or Both")
}

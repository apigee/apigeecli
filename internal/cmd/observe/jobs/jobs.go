// Copyright 2024 Google LLC
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

package jobs

import "github.com/spf13/cobra"

// ObservationJobCmd to manage apis
var ObservationJobCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Manage Observation jobs for Shadow API Discovery",
	Long:  "Manage Observation jobs for Shadow API Discovery",
}

var org, region string

func init() {
	ObservationJobCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	ObservationJobCmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "API Observation region name")

	ObservationJobCmd.AddCommand(CrtCmd)
	ObservationJobCmd.AddCommand(GetCmd)
	ObservationJobCmd.AddCommand(DeleteCmd)
	ObservationJobCmd.AddCommand(ListCmd)
	ObservationJobCmd.AddCommand(EnableCmd)
	ObservationJobCmd.AddCommand(DisableCmd)
}

// Copyright 2022 Google LLC
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

package org

import (
	"github.com/spf13/cobra"
)

// ReportCmd to manage org reprots
var ReportCmd = &cobra.Command{
	Use:     "reports",
	Aliases: []string{"orgs"},
	Short:   "Report Apigee Org Usage",
	Long:    "Report Apigee Org Usage",
}

const (
	apiCallsHeader  = "\tAPI CALLS"
	proxyTypeHeader = "\tEXTENSIBLE PROXY\tSTANDARD PROXY"
)

func init() {
	ReportCmd.AddCommand(MonthlyCmd)
	ReportCmd.AddCommand(YearlyCmd)
}

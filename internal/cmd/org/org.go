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

package org

import (
	"github.com/spf13/cobra"
)

// Cmd to manage orgs
var Cmd = &cobra.Command{
	Use:     "organizations",
	Aliases: []string{"orgs"},
	Short:   "Manage Apigee Orgs",
	Long:    "Manage Apigee Orgs",
}

var org, region string

func init() {
	Cmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "Apigee control plane region name; default is https://apigee.googleapis.com")

	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(MartCmd)
	Cmd.AddCommand(AcCmd)
	Cmd.AddCommand(PropCmd)
	Cmd.AddCommand(ObCmd)
	Cmd.AddCommand(IngressCmd)
	Cmd.AddCommand(ExportCmd)
	Cmd.AddCommand(ImportCmd)
	Cmd.AddCommand(UpdateCmd)
	Cmd.AddCommand(SetAddonCmd)
	Cmd.AddCommand(ReportCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(DeployCmd)
}

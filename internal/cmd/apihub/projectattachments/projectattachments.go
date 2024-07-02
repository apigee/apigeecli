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

package projectattachments

import (
	"github.com/spf13/cobra"
)

// ProjectAttachmentCmd to manage apis
var ProjectAttachmentCmd = &cobra.Command{
	Use:   "project-attachments",
	Short: "Manage Runtime Project Attachments to API Hub",
	Long:  "Manage Runtime Project Attachments to API Hub",
}

var org, region string

func init() {
	ProjectAttachmentCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	ProjectAttachmentCmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "API Hub region name")

	ProjectAttachmentCmd.AddCommand(CrtCmd)
	ProjectAttachmentCmd.AddCommand(GetCmd)
	ProjectAttachmentCmd.AddCommand(DelCmd)
	ProjectAttachmentCmd.AddCommand(ListCmd)

	_ = ProjectAttachmentCmd.MarkFlagRequired("org")
	_ = ProjectAttachmentCmd.MarkFlagRequired("region")
}

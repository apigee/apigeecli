// Copyright 2021 Google LLC
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

package eptattachment

import (
	"fmt"
	"regexp"

	"internal/apiclient"

	"internal/client/eptattachment"

	"github.com/spf13/cobra"
)

// CreateCmd to list endpoint attachments
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new service endpoint",
	Long:  "Create a new service endpoint in Apigee for a service attachment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		re := regexp.MustCompile(`projects\/([a-zA-Z0-9_-]+)\/regions` +
			`\/([a-zA-Z0-9_-]+)\/serviceAttachments\/([a-zA-Z0-9_-]+)`)
		ok := re.Match([]byte(location))
		if !ok {
			return fmt.Errorf("disk encryption key must be of the format " +
				"projects/{project-id}/regions/{location}/serviceAttachments/{sa-name}")
		}
		_, err = eptattachment.Create(name, serviceAttachment, location)
		return err
	},
}

var location, serviceAttachment string

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the service endpoint")
	CreateCmd.Flags().StringVarP(&location, "location", "l",
		"", "Location of the service endpoint")
	CreateCmd.Flags().StringVarP(&serviceAttachment, "service-attachment", "s",
		"", "Service attachment url: projects/{project-id}/regions/{location}/serviceAttachments/{sa-name}")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("service-attachment")
	_ = CreateCmd.MarkFlagRequired("location")
}

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
	"fmt"
	"internal/apiclient"
	"internal/client/developers"

	"github.com/spf13/cobra"
)

// UpdateCmd to update developer
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an Apigee developer configuration",
	Long:  "Update an Apigee developer configuration",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = developers.Update(email, firstName, lastName, userName, cmd.Flag("status").Value.String(), attrs)
		return
	},
}

type developerStatus string

var status developerStatus

const (
	active   developerStatus = "active"
	inactive developerStatus = "inactive"
)

func (s *developerStatus) String() string {
	return string(*s)
}

func (s *developerStatus) Set(r string) error {
	switch r {
	case "active", "inactive":
		*s = developerStatus(r)
		return nil
	default:
		return fmt.Errorf("must be %s or %s", active, inactive)
	}
}

func (s *developerStatus) Type() string {
	return "developerStatus"
}

func init() {
	UpdateCmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")
	UpdateCmd.Flags().StringVarP(&firstName, "first", "f",
		"", "The first name of the developer")
	UpdateCmd.Flags().StringVarP(&lastName, "last", "s",
		"", "The last name of the developer")
	UpdateCmd.Flags().StringVarP(&userName, "user", "u",
		"", "The username of the developer")
	UpdateCmd.Flags().Var(&status, "status",
		fmt.Sprintf("must be %s or %s", active, inactive))
	UpdateCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = UpdateCmd.MarkFlagRequired("email")
}

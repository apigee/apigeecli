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

package spaces

import (
	"github.com/spf13/cobra"
)

// IamCmd to manage manage env iam permissions
var IamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM permissions for a space",
	Long:  "Manage IAM permissions for a space",
}

var memberName, role, memberType, space string

func init() {
	IamCmd.PersistentFlags().StringVarP(&space, "space", "",
		"", "Apigee space name")

	_ = IamCmd.MarkPersistentFlagRequired("space")

	IamCmd.AddCommand(GetIamCmd)
	IamCmd.AddCommand(RemoveRoleCmd)
	IamCmd.AddCommand(SetEditorCmd)
	IamCmd.AddCommand(TestIamCmd)
}

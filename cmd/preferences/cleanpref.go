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

package preferences

import (
	"internal/apiclient"

	"github.com/spf13/cobra"
)

// Cmd to get org details
var CleanCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove preferences file from the home dir",
	Long:  "Remove preferences file from the home dir",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.DeletePreferencesFile()
	},
}

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

package fetchsf

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/client/sharedflows"
)

//Cmd to download shared flow
var Cmd = &cobra.Command{
	Use:   "fetch",
	Short: "Returns a zip-formatted shared flow bundle ",
	Long:  "Returns a zip-formatted shared flow bundle of code and config files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sharedflows.Fetch(name, revision)
	},
}

var name, revision string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Shared flow Bundle Name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "Shared flow revision")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("rev")
}

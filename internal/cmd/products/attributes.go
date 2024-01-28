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

package products

import (
	"github.com/spf13/cobra"
)

// AttributesCmd to manage tracing of apis
var AttributesCmd = &cobra.Command{
	Use:     "attributes",
	Aliases: []string{"attrs"},
	Short:   "Manage API Product Attributes",
	Long:    "Manage API Product Attributes",
}

var attrName string

func init() {
	AttributesCmd.PersistentFlags().StringVarP(&name, "name", "n",
		"", "API Product name")

	_ = AttributesCmd.MarkPersistentFlagRequired("name")

	AttributesCmd.AddCommand(DelAttrCmd)
	AttributesCmd.AddCommand(ListAttrCmd)
	AttributesCmd.AddCommand(GetAttrCmd)
	AttributesCmd.AddCommand(UpdAttrCmd)
}

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

package kvm

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// EntryCmd to manage kvm entries
var EntryCmd = &cobra.Command{
	Use:   "entries",
	Short: "Manage Key Value Map Entries",
	Long:  "Manage Key Value Map Entries",
}

var mapName, keyName string

func init() {
	EntryCmd.AddCommand(CreateEntryCmd)
	EntryCmd.AddCommand(GetEntryCmd)
	EntryCmd.AddCommand(DelEntryCmd)
	EntryCmd.AddCommand(ListEntryCmd)
	EntryCmd.AddCommand(ExpEntryCmd)
	EntryCmd.AddCommand(ImpEntryCmd)
	EntryCmd.AddCommand(UpdateEntryCmd)
}

func getKVMString(value string) string {
	var err error
	// convert any boolean, float or integer to string
	if _, err = strconv.ParseBool(value); err == nil {
		value = fmt.Sprintf("\"%s\"", value)
	} else if _, err = strconv.ParseInt(value, 10, 0); err == nil {
		value = fmt.Sprintf("\"%s\"", value)
	} else if _, err = strconv.ParseFloat(value, 0); err == nil {
		value = fmt.Sprintf("\"%s\"", value)
	}
	return value
}

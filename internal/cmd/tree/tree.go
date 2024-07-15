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

package tree

import (
	"internal/clilog"

	"github.com/spf13/cobra"
)

// Cmd to to print all the commands in apigeecli
var Cmd = &cobra.Command{
	Use:   "tree",
	Short: "Prints apigeecli command Tree",
	Long:  "Prints apigeecli command Tree",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		printCommand(cmd.Parent().Root(), 0, true)
		return
	},
}

func printCommand(cmd *cobra.Command, depth int, last bool) {
	var prefix string
	if depth > 0 {
		if last {
			prefix = "└── "
		} else {
			prefix = "├── "
		}

		for i := 0; i < depth-1; i++ {
			prefix = (string(rune(0x2502)) + "   ") + prefix // Vertical bar
		}
	}

	clilog.Info.Printf("%s%s\n", prefix, cmd.Use+" - "+cmd.Short)

	for i, subCmd := range cmd.Commands() {
		isLast := i == len(cmd.Commands())-1
		printCommand(subCmd, depth+1, isLast)
	}
}

package tree

import (
	"internal/clilog"

	"github.com/spf13/cobra"
)

// Cmd to manage token to interact with apigee.googleapis.com
var Cmd = &cobra.Command{
	Use:   "tree",
	Short: "Print Tree",
	Long:  "Print Tree",
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

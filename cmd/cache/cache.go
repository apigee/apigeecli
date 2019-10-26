package cache

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/cache/delcache"
	"github.com/srinandan/apigeecli/cmd/cache/listcache"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage caches within an Apigee environment",
	Long:  "Manage caches within an Apigee environment",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(listcache.Cmd)
	Cmd.AddCommand(delcache.Cmd)
}

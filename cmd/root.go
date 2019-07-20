package cmd

import (
	"./apis"
	"./apps"
	"./developers"
	"./env"
	flowhooks "./flowhooks"
	"./kvm"
	"./org"
	"./products"
	"./shared"
	"./sharedflows"
	"./sync"
	targetservers "./targetservers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:     "apigeeapi",
	Version: "0.1",
	Short:   "Utility to work with Apigee APIs.",
	Long:    "This command lets you interact with Apigee APIs.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		shared.Init()
		return shared.SetAccessToken()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&shared.LogInfo, "log", "l",
		false, "Log Information")

	RootCmd.PersistentFlags().StringVarP(&shared.RootArgs.Token, "token", "t",
		"", "Google OAuth Token")
	_ = viper.BindPFlag("token", RootCmd.PersistentFlags().Lookup("token"))

	RootCmd.PersistentFlags().StringVarP(&shared.RootArgs.ServiceAccount, "account", "a",
		"", "Path Service Account private key in JSON")
	_ = viper.BindPFlag("account", RootCmd.PersistentFlags().Lookup("account"))

	RootCmd.AddCommand(apis.Cmd)
	RootCmd.AddCommand(org.Cmd)
	RootCmd.AddCommand(sync.Cmd)
	RootCmd.AddCommand(env.Cmd)
	RootCmd.AddCommand(products.Cmd)
	RootCmd.AddCommand(developers.Cmd)
	RootCmd.AddCommand(apps.Cmd)
	RootCmd.AddCommand(sharedflows.Cmd)
	RootCmd.AddCommand(kvm.Cmd)
	RootCmd.AddCommand(flowhooks.Cmd)
	RootCmd.AddCommand(targetservers.Cmd)
}

func initConfig() {
	viper.SetEnvPrefix("APIGEE")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetConfigType("json")
}

// GetRootCmd returns the root of the cobra command-tree.
func GetRootCmd() *cobra.Command {
	return RootCmd
}

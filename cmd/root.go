package cmd

import (
	"./apis"
	"./apps"
	"./developers"
	"./env"
	"./org"
	"./products"
	"./shared"
	"./sharedflows"
	"./sync"
	"fmt"
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

		if shared.RootArgs.Token == "" && shared.RootArgs.ServiceAccount == "" {
			return fmt.Errorf("Either token or service account must be provided")
		} else {
			if shared.RootArgs.ServiceAccount != "" {
				viper.SetConfigFile(shared.RootArgs.ServiceAccount)
				err := viper.ReadInConfig() // Find and read the config file
				if err != nil {             // Handle errors reading the config file
					return fmt.Errorf("Fatal error config file: %s \n", err)
				} else {
					if viper.Get("private_key") == "" {
						return fmt.Errorf("Fatal error: Private key missing in the service account")
					}
					if viper.Get("client_email") == "" {
						return fmt.Errorf("Fatal error: client email missing in the service account")
					}
					_, err = shared.GenerateAccessToken()
					if err != nil {
						return fmt.Errorf("Fatal error generating access token: %s \n", err)
					} else {
						return nil
					}
				}
			}
			return nil
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&shared.LogInfo, "log", "l",
		false, "Log Information")

	RootCmd.PersistentFlags().StringVarP(&shared.RootArgs.Token, "token", "t",
		"", "Google OAuth Token")
	viper.BindPFlag("token", RootCmd.PersistentFlags().Lookup("token"))

	RootCmd.PersistentFlags().StringVarP(&shared.RootArgs.ServiceAccount, "account", "a",
		"", "Path Service Account private key in JSON")
	viper.BindPFlag("account", RootCmd.PersistentFlags().Lookup("account"))

	RootCmd.AddCommand(apis.Cmd)
	RootCmd.AddCommand(org.Cmd)
	RootCmd.AddCommand(sync.Cmd)
	RootCmd.AddCommand(env.Cmd)
	RootCmd.AddCommand(products.Cmd)
	RootCmd.AddCommand(developers.Cmd)
	RootCmd.AddCommand(apps.Cmd)
	RootCmd.AddCommand(sharedflows.Cmd)
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

package org

import (
	"./getorg"
	"./listorgs"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "orgs",
	Short: "Manage Apigee Orgs",
	Long:  "Manage Apigee Orgs",
}

var expand = false
var count string

func init() {

	Cmd.AddCommand(listorgs.Cmd)
	Cmd.AddCommand(getorg.Cmd)
}

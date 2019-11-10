package org

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/org/createorg"
	"github.com/srinandan/apigeecli/cmd/org/getorg"
	"github.com/srinandan/apigeecli/cmd/org/listorgs"
	"github.com/srinandan/apigeecli/cmd/org/setmart"
)

//Cmd to manage orgs
var Cmd = &cobra.Command{
	Use:   "orgs",
	Short: "Manage Apigee Orgs",
	Long:  "Manage Apigee Orgs",
}

func init() {
	Cmd.AddCommand(createorg.Cmd)
	Cmd.AddCommand(listorgs.Cmd)
	Cmd.AddCommand(getorg.Cmd)
	Cmd.AddCommand(setmart.Cmd)
}

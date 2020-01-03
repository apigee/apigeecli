package iam

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/iam/createall"
	"github.com/srinandan/apigeecli/cmd/iam/createax"
	"github.com/srinandan/apigeecli/cmd/iam/createcass"
	"github.com/srinandan/apigeecli/cmd/iam/createconnect"
	"github.com/srinandan/apigeecli/cmd/iam/createlogger"
	"github.com/srinandan/apigeecli/cmd/iam/createmart"
	"github.com/srinandan/apigeecli/cmd/iam/createmetrics"
	"github.com/srinandan/apigeecli/cmd/iam/createsync"
)

//Cmd to manage orgs
var Cmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM permissions for Apigee",
	Long: "Manage IAM permissions for Apigee. The SA to run this command requires Security Admin, " +
		"Create Service Accounts and Service Account Key Admin roles",
}

func init() {
	Cmd.AddCommand(createall.Cmd)
	Cmd.AddCommand(createax.Cmd)
	Cmd.AddCommand(createcass.Cmd)
	Cmd.AddCommand(createconnect.Cmd)
	Cmd.AddCommand(createlogger.Cmd)
	Cmd.AddCommand(createmart.Cmd)
	Cmd.AddCommand(createmetrics.Cmd)
	Cmd.AddCommand(createsync.Cmd)
}

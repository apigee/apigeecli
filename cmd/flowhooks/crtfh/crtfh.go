package crtfh

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
	"strconv"
	"strings"
)

var Cmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a flowhook",
	Long:  "Attach a flowhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)

		flowhook := []string{}

		flowhook = append(flowhook, "\"name\":\""+name+"\"")

		if description != "" {
			flowhook = append(flowhook, "\"description\":\""+description+"\"")
		}

		flowhook = append(flowhook, "\"sharedFlow\":\""+sharedflow+"\"")

		if continueOnErr {
			flowhook = append(flowhook, "\"continueOnError\":"+strconv.FormatBool(continueOnErr))
		}

		payload := "{" + strings.Join(flowhook, ",") + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "flowhooks", name)
		return shared.HttpClient(u.String(), payload, "PUT")
	},
}

var name, description, sharedflow string
var continueOnErr bool

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Flowhook point")
	Cmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the flowhook")
	Cmd.Flags().StringVarP(&sharedflow, "sharedflow", "s",
		"", "Sharedflow name")
	Cmd.Flags().BoolVarP(&continueOnErr, "continue", "c",
		true, "Continue on error")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("sharedflow")
}

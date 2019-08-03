package crtts

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
	"strconv"
	"strings"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Target Server",
	Long:  "Create a Target Server",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)

		targetserver := []string{}

		targetserver = append(targetserver, "\"name\":\""+name+"\"")

		if description != "" {
			targetserver = append(targetserver, "\"description\":\""+description+"\"")
		}

		targetserver = append(targetserver, "\"host\":\""+host+"\"")
		targetserver = append(targetserver, "\"port\":"+strconv.Itoa(port))

		if enable {
			targetserver = append(targetserver, "\"isEnabled\":"+strconv.FormatBool(enable))
		}

		payload := "{" + strings.Join(targetserver, ",") + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "targetservers")
		_, err = shared.HttpClient(true, u.String(), payload) 
		return 

	},
}

var name, description, host string
var enable bool
var port int

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the targetserver")
	Cmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the Target Server")
	Cmd.Flags().StringVarP(&host, "host", "s",
		"", "Host name of the target")
	Cmd.Flags().BoolVarP(&enable, "enable", "b",
		true, "Enabling/disabling a TargetServer")
	Cmd.Flags().IntVarP(&port, "port", "p",
		80, "port number")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("host")
}

package crtts

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
	"strconv"
	"strings"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Target Server",
	Long:  "Create a Target Server",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)

		targetserver := []string{}

		targetserver = append(targetserver, "\"name\":\""+name+"\"")

		if description != "" {
			targetserver = append(targetserver, "\"description\":\""+description+"\"")
		}

		targetserver = append(targetserver, "\"host\":\""+host+"\"")
		targetserver = append(targetserver, "\"port\":"+strconv.Itoa(port))

		if enable  {
			targetserver = append(targetserver, "\"isEnabled\":"+strconv.FormatBool(enable))
		}

		payload := "{" + strings.Join(targetserver, ",") + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "targetservers")
		shared.HttpClient(u.String(), payload)
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

	Cmd.MarkFlagRequired("name")
	Cmd.MarkFlagRequired("host")
}


package getdev

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/thedevsaddam/gojsonq"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get App in an Organization by App ID",
	Long:  "Returns the app profile for the specified app ID.",
	Args: func(cmd *cobra.Command, args []string) error {
		if appId == "" && name == "" {
			return fmt.Errorf("Pass either name or appId")
		}
		if appId != "" && name != "" {
			return fmt.Errorf("name and appId cannot be used together")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		if appId != "" {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "apps", appId)
			return shared.HttpClient(u.String())
		} else {
			//search by name is not implement; use list and return the appropriate app
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "apps")
			q := u.Query()
			q.Set("expand", "true")
			q.Set("includeCred", "false")
			u.RawQuery = q.Encode()
			//don't print the list
			shared.RootArgs.Print = false
			err := shared.HttpClient(u.String())
			if err != nil {
				return err
			}
			jq := gojsonq.New().JSONString(string(shared.RootArgs.Body)).From("app").Where("name", "eq", name)
			out := jq.Get()
			outBytes, err := json.Marshal(out)
			if err != nil {
				return err
			}
			//print the item
			shared.RootArgs.Print = true
			return shared.PrettyPrint(outBytes)
		}
	},
}

var appId, name string

func init() {

	Cmd.Flags().StringVarP(&appId, "appId", "i",
		"", "Developer app id")
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Developer app name")
}

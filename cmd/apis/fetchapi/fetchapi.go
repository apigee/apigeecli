package fetchapi

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "fetch",
	Short: "Returns a zip-formatted proxy bundle ",
	Long:  "Returns a zip-formatted proxy bundle of code and config files",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		q := u.Query()
		q.Set("format", "bundle")
		u.RawQuery = q.Encode()
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis", name, "revisions", revision)
		_ = shared.DownloadResource(u.String(), name)
	},
}

var name, revision string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy Bundle Name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "API Proxy revision")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("revision")
}

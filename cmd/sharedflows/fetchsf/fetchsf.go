package fetchsf

import (
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "fetch",
	Short: "Returns a zip-formatted shared flow bundle ",
	Long:  "Returns a zip-formatted shared flow bundle of code and config files",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		q := u.Query()
		q.Set("format", "bundle")
		u.RawQuery = q.Encode()
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "sharedflows", name, "revisions", revision)
		return shared.DownloadResource(u.String(), name)
	},
}

var name, revision string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Shared flow Bundle Name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "Shared flow revision")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("rev")
}

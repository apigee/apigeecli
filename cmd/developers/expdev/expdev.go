package expdev

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export Developers to a file",
	Long:  "Export Developers to a file",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		const exportFileName = "developers.json"

		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers")

		q := u.Query()
		q.Set("expand", "true")

		u.RawQuery = q.Encode()
		//don't print to sysout
		respBody, err := shared.HttpClient(false, u.String())
		if err != nil {
			return err
		}

		return shared.WriteByteArrayToFile(exportFileName, false, respBody)
	},
}

func init() {

}

package crtkvm

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
	Short: "Create an environment scoped KVM Map",
	Long:  "Create an environment scoped KVM Map",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		kvm := []string{}

		kvm = append(kvm, "\"name\":\""+name+"\"")
		if encrypt {
			kvm = append(kvm, "\"encrypted\":"+strconv.FormatBool(encrypt))
		}
		payload := "{" + strings.Join(kvm, ",") + "}"
		shared.Info.Println(payload)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keyvaluemaps")
		_, err = shared.HttpClient(true, u.String(), payload) 
		return 		
	},
}

var name string
var encrypt bool

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Environment name")
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "KVM Map name")
	Cmd.Flags().BoolVarP(&encrypt, "encrypt", "c",
		false, "Enable cncrypted KVM")

	_ = Cmd.MarkFlagRequired("env")
	_ = Cmd.MarkFlagRequired("name")
}

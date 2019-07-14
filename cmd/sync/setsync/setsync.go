package setsync

import (
	"fmt"
	"net/url"
	"../../shared"
	"path"
	"encoding/json"
	"github.com/spf13/cobra"
	"strings"
)

//{"identities":["serviceAccount:srinandans-apigee@srinandans-apigee.iam.gserviceaccount.com"]}

type IAMIdentities struct {
	Identities []string `json:"identities,omitempty"`
}

var identity string

var Cmd = &cobra.Command{
	Use:   "set",
	Short: "Set identity with access to control plane resources",
	Long:  "Set identity with access to control plane resources",
	Args: func(cmd *cobra.Command, args []string) error {
		if strings.Contains(identity, ".iam.gserviceaccount.com") {
			return nil
		} else {
			return fmt.Errorf("identity[0] must have .iam.gserviceaccount.com suffix and should not be a Google managed service account: %s", identity)
		}
	}, 
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org+":setSyncAuthorization")
		identity = validate(identity)
		identities := IAMIdentities{}
		identities.Identities = append(identities.Identities, identity)
		payload, _ := json.Marshal(&identities)
		fmt.Println(string(payload))
		shared.PostHttpClient(u.String(), shared.RootArgs.Token, string(payload))
	},
}

func init() {


	Cmd.Flags().StringVarP(&identity, "ity", "i",
	"", "IAM Identity")

	Cmd.MarkFlagRequired("ity")		
}

func validate(i string) string {
	if strings.Contains(i, "serviceAccount:")  {
		return i
	} else {
		return "serviceAccount:"+i
	}
}
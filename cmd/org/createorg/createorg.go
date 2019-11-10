package createorg

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to get org details
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Apigee Org",
	Long:  "Create a new Apigee Org; Your GCP project must be whitelist for this operation",
	Args: func(cmd *cobra.Command, args []string) error {
		if !validRegion(region) {
			return fmt.Errorf("Invalid analytics region."+
				" Analytics region must be one of : %v", analyticsRegions)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(baseURL)
		u.Path = path.Join(u.Path)
		q := u.Query()
		q.Set("parent", "projects/"+shared.RootArgs.ProjectID)
		u.RawQuery = q.Encode()

		orgPayload := []string{}
		orgPayload = append(orgPayload, "\"name\":\""+shared.RootArgs.Org+"\"")
		orgPayload = append(orgPayload, "\"analyticsRegion\":\""+region+"\"")

		payload := "{" + strings.Join(orgPayload, ",") + "}"
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}

var analyticsRegions = [...]string{"asia-east1", "asia-east1", "asia-northeast1", "asia-southeast1",
	"europe-west1", "us-central1", "us-east1", "us-east4", "us-west1", "australia-southeast1",
	"europe-west2"}

func validRegion(region string) bool {
	for _, r := range analyticsRegions {
		if region == r {
			return true
		}
	}
	return false
}

const baseURL = "https://apigee.googleapis.com/v1/organizations"

var region string

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.Flags().StringVarP(&region, "reg", "r",
		"", "Analytics region name")
	Cmd.Flags().StringVarP(&shared.RootArgs.ProjectID, "prj", "p",
		"", "GCP Project ID")

	_ = Cmd.MarkFlagRequired("prj")
	_ = Cmd.MarkFlagRequired("org")
	_ = Cmd.MarkFlagRequired("reg")
}

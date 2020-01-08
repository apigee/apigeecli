package main

import (
	"fmt"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/orgs"
)

func main() {
	args := apiclient.ApigeeClientOptions{
		Org:            "gcp-project-id",
		ServiceAccount: "path-to-service-account.json",
		SkipLogInfo:    true,
	}
	apiclient.NewApigeeClient(args)

	err := apiclient.SetAccessToken()
	if err != nil {
		fmt.Println(err)
		return
	}

	respBody, err := orgs.List()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(respBody))
}

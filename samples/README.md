# apigeecli Client Sample

apigeecli can be used as a golang based client library. Here is a sample to get a list of orgs.

```go

package main

import (
	"fmt"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/orgs"
)

func main() {

	//set client options
	apiclient.NewApigeeClient(apiclient.ApigeeClientOptions{
		Org:            "apigee-org-name",
		ServiceAccount: "path-to-service-account.json",
		SkipLogInfo:    true,                             //skip printing client logs
	})

	//invoke list of orgs
	respBody, err := orgs.List()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(respBody))
}
```
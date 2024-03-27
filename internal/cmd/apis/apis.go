// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apis

import (
	"github.com/spf13/cobra"
)

// Cmd to manage apis
var Cmd = &cobra.Command{
	Use:   "apis",
	Short: "Manage Apigee API proxies in an org",
	Long:  "Manage Apigee API proxies in an org",
}

var (
	org, region, env, name string
	conn, revision         int
)

const zipExt = ".zip"

var examples = []string{"apigeecli apis import -f samples/apis  --default-token",
	`apigeecli apis create oas -n petstore \
--oas-base-folderpath=./samples \
--oas-name=petstore.yaml --add-cors=true \
--google-idtoken-aud-literal=https://sample.run.app  --default-token`,
	`apigeecli apis create oas -n petstore \
-f ./samples/petstore.yaml \
--add-cors=true --env=$env --wait=true --default-token`,
	"apigeecli apis create graphql -n petstore --action=parse_verify --add-cors=true  --default-token",
	"apigeecli apis create integration -n sample --trigger=sample  --default-token",
	`apigeecli apis create github --name=petstore \
--owner=apigee \
--repo=api-platform-samples \
--proxy-path=sample-proxies/apikey  --default-token`}

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "Apigee control plane region name; default is https://apigee.googleapis.com")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(ListDepCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(ExpCmd)
	Cmd.AddCommand(DepCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(FetCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(ImpCmd)
	Cmd.AddCommand(UndepCmd)
	Cmd.AddCommand(TraceCmd)
	Cmd.AddCommand(CleanCmd)
	Cmd.AddCommand(KvmCmd)
	Cmd.AddCommand(UpdateCmd)
	Cmd.AddCommand(CloneCmd)
}

func GetExample(i int) string {
	return examples[i]
}

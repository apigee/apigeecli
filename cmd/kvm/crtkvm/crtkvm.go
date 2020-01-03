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

package crtkvm

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
	"strconv"
	"strings"
)

//Cmd to create kvms
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

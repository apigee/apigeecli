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

package crtfh

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to create flow hooks
var Cmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a flowhook",
	Long:  "Attach a flowhook",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)

		flowhook := []string{}

		flowhook = append(flowhook, "\"name\":\""+name+"\"")

		if description != "" {
			flowhook = append(flowhook, "\"description\":\""+description+"\"")
		}

		flowhook = append(flowhook, "\"sharedFlow\":\""+sharedflow+"\"")

		if continueOnErr {
			flowhook = append(flowhook, "\"continueOnError\":"+strconv.FormatBool(continueOnErr))
		}

		payload := "{" + strings.Join(flowhook, ",") + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "flowhooks", name)
		_, err = shared.HttpClient(true, u.String(), payload, "PUT")
		return
	},
}

var name, description, sharedflow string
var continueOnErr bool

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Flowhook point")
	Cmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the flowhook")
	Cmd.Flags().StringVarP(&sharedflow, "sharedflow", "s",
		"", "Sharedflow name")
	Cmd.Flags().BoolVarP(&continueOnErr, "continue", "c",
		true, "Continue on error")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("sharedflow")
}

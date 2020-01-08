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
	"github.com/srinandan/apigeecli/apiclient"
	crtapi "github.com/srinandan/apigeecli/cmd/apis/crtapi"
	delapi "github.com/srinandan/apigeecli/cmd/apis/delapi"
	"github.com/srinandan/apigeecli/cmd/apis/depapi"
	"github.com/srinandan/apigeecli/cmd/apis/expapis"
	fetch "github.com/srinandan/apigeecli/cmd/apis/fetchapi"
	"github.com/srinandan/apigeecli/cmd/apis/getapi"
	impapis "github.com/srinandan/apigeecli/cmd/apis/impapis"
	"github.com/srinandan/apigeecli/cmd/apis/listapis"
	"github.com/srinandan/apigeecli/cmd/apis/listdeploy"
	traceapi "github.com/srinandan/apigeecli/cmd/apis/traceapi"
	"github.com/srinandan/apigeecli/cmd/apis/undepapi"
)

//Cmd to manage apis
var Cmd = &cobra.Command{
	Use:   "apis",
	Short: "Manage Apigee API proxies in an org",
	Long:  "Manage Apigee API proxies in an org",
}

func init() {

	Cmd.PersistentFlags().StringVarP(apiclient.GetApigeeOrgP(), "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listapis.Cmd)
	Cmd.AddCommand(listdeploy.Cmd)
	Cmd.AddCommand(crtapi.Cmd)
	Cmd.AddCommand(expapis.Cmd)
	Cmd.AddCommand(depapi.Cmd)
	Cmd.AddCommand(delapi.Cmd)
	Cmd.AddCommand(fetch.Cmd)
	Cmd.AddCommand(getapi.Cmd)
	Cmd.AddCommand(impapis.Cmd)
	Cmd.AddCommand(undepapi.Cmd)
	Cmd.AddCommand(traceapi.Cmd)
}

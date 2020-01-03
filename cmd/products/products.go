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

package products

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/products/crtprod"
	"github.com/srinandan/apigeecli/cmd/products/delprod"
	"github.com/srinandan/apigeecli/cmd/products/expprod"
	"github.com/srinandan/apigeecli/cmd/products/getprod"
	"github.com/srinandan/apigeecli/cmd/products/impprod"
	"github.com/srinandan/apigeecli/cmd/products/listproducts"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage products
var Cmd = &cobra.Command{
	Use:     "products",
	Aliases: []string{"prods"},
	Short:   "Manage Apigee API products",
	Long:    "Manage Apigee API products",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listproducts.Cmd)
	Cmd.AddCommand(getprod.Cmd)
	Cmd.AddCommand(delprod.Cmd)
	Cmd.AddCommand(crtprod.Cmd)
	Cmd.AddCommand(impprod.Cmd)
	Cmd.AddCommand(expprod.Cmd)
}

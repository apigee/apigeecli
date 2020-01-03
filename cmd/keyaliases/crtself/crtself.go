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

package crtself

import (
	"encoding/json"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to create key aliases
var Cmd = &cobra.Command{
	Use:   "create-self-signed",
	Short: "Create a Key Alias from self-seigned cert",
	Long:  "Create a Key Alias by generating a self-signed cert",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env,
			"keystores", shared.RootArgs.AliasName, "aliases")

		var jsonPayload map[string]interface{}
		err = json.Unmarshal([]byte(payload), &jsonPayload)

		if err != nil {
			return
		}

		q := u.Query()
		q.Set("format", certFormat)
		q.Set("alias", aliasName)

		if ignoreNewLine {
			q.Set("ignoreNewlineValidation", "true")
		}
		if ignoreExpiry {
			q.Set("ignoreExpiryValidation", "true")
		}
		u.RawQuery = q.Encode()
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}

const certFormat = "selfsignedcert"

var aliasName, payload string
var ignoreNewLine, ignoreExpiry bool

func init() {

	Cmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key alias")
	Cmd.Flags().StringVarP(&payload, "cert", "c",
		"", "Certificate in JSON format")
	Cmd.Flags().BoolVarP(&ignoreExpiry, "exp", "x",
		false, "Ignore expiry validation")
	Cmd.Flags().BoolVarP(&ignoreNewLine, "nl", "w",
		false, "Ignore new line in cert chain")

	_ = Cmd.MarkFlagRequired("alias")
	_ = Cmd.MarkFlagRequired("cert")
}

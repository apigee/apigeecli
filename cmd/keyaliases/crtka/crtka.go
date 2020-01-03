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

package crtka

import (
	"fmt"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to create key aliases
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Alias from PEM, JAR or PKCS12 file",
	Long:  "Create a Key Alias from PEM, JAR or PKCS12 file",
	Args: func(cmd *cobra.Command, args []string) error {
		if !validate() {
			return fmt.Errorf("certificate format must be one of keycertfile, pkcs12 or keycertjar")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env,
			"keystores", shared.RootArgs.AliasName, "aliases")

		q := u.Query()
		q.Set("format", format)
		q.Set("alias", aliasName)

		if ignoreNewLine {
			q.Set("ignoreNewlineValidation", "true")
		}
		if ignoreExpiry {
			q.Set("ignoreExpiryValidation", "true")
		}
		if password != "" {
			q.Set("password", password)
		}

		u.RawQuery = q.Encode()

		if format == "keycertfile" {
			_, err = shared.PostHttpOctet(true, u.String(), aliasName+".pem")
		} else if format == "pkcs12" {
			_, err = shared.PostHttpOctet(true, u.String(), aliasName+".pfx")
		}
		return
	},
}

var certFormats = [2]string{"keycertfile", "pkcs12"}

var aliasName, format, password string
var ignoreNewLine, ignoreExpiry bool

func init() {

	Cmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key alias")
	Cmd.Flags().StringVarP(&format, "format", "f",
		"", "Format of the certificate")
	Cmd.Flags().StringVarP(&password, "password", "p",
		"", "PKCS12 password")
	Cmd.Flags().BoolVarP(&ignoreExpiry, "exp", "x",
		false, "Ignore expiry validation")
	Cmd.Flags().BoolVarP(&ignoreNewLine, "nl", "w",
		false, "Ignore new line in cert chain")

	_ = Cmd.MarkFlagRequired("alias")
	_ = Cmd.MarkFlagRequired("format")
}

func validate() bool {
	for _, frmt := range certFormats {
		if format == frmt {
			return true
		}
	}
	return false
}

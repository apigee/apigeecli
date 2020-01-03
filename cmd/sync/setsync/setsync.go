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

package setsync

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//{"identities":["serviceAccount:srinandans-apigee@srinandans-apigee.iam.gserviceaccount.com"]}

type iAMIdentities struct {
	Identities []string `json:"identities,omitempty"`
}

var identity string

//Cmd to set identities
var Cmd = &cobra.Command{
	Use:   "set",
	Short: "Set identity with access to control plane resources",
	Long:  "Set identity with access to control plane resources",
	Args: func(cmd *cobra.Command, args []string) error {
		if !strings.Contains(identity, ".iam.gserviceaccount.com") {
			return fmt.Errorf("identity[0] must have .iam.gserviceaccount.com suffix"+
				" and should not be a Google managed service account: %s", identity)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org+":setSyncAuthorization")
		identity = validate(identity)
		identities := iAMIdentities{}
		identities.Identities = append(identities.Identities, identity)
		payload, _ := json.Marshal(&identities)
		_, err = shared.HttpClient(true, u.String(), string(payload))
		return

	},
}

func init() {

	Cmd.Flags().StringVarP(&identity, "ity", "i",
		"", "IAM Identity")

	_ = Cmd.MarkFlagRequired("ity")
}

func validate(i string) string {
	if strings.Contains(i, "serviceAccount:") {
		return i
	}
	return "serviceAccount:" + i
}

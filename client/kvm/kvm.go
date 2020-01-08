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

package kvm

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/srinandan/apigeecli/apiclient"
)

func Create(name string, encrypt bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	kvm := []string{}

	kvm = append(kvm, "\"name\":\""+name+"\"")
	if encrypt {
		kvm = append(kvm, "\"encrypted\":"+strconv.FormatBool(encrypt))
	}
	payload := "{" + strings.Join(kvm, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps")
	respBody, err = apiclient.HttpClient(true, u.String(), payload)
	return respBody, err
}

func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", name)
	respBody, err = apiclient.HttpClient(true, u.String(), "", "DELETE")
	return respBody, err
}

func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps")
	respBody, err = apiclient.HttpClient(true, u.String())
	return respBody, err
}

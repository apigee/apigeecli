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

package flowhooks

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/srinandan/apigeecli/apiclient"
)

func Attach(name string, description string, sharedflow string, continueOnErr bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

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
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "flowhooks", name)
	respBody, err = apiclient.HttpClient(true, u.String(), payload, "PUT")
	return respBody, err
}

func Detach(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "flowhooks", name)
	respBody, err = apiclient.HttpClient(true, u.String(), "", "DELETE")
	return respBody, err
}

func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "flowhooks", name)
	respBody, err = apiclient.HttpClient(true, u.String())
	return respBody, err
}

func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "flowhooks")
	respBody, err = apiclient.HttpClient(true, u.String())
	return respBody, err
}

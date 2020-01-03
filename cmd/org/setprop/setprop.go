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

package setprop

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/srinandan/apigeecli/cmd/types"
)

func SetOrgProperty(name string, value string) (err error) {
	u, _ := url.Parse(shared.BaseURL)
	u.Path = path.Join(u.Path, shared.RootArgs.Org)
	//get org details
	orgBody, err := shared.HttpClient(false, u.String())
	if err != nil {
		return err
	}

	org := types.Org{}
	err = json.Unmarshal(orgBody, &org)
	if err != nil {
		return err
	}

	//check if the property exists
	found := false
	for i, properties := range org.Properties.Property {
		if properties.Name == name {
			fmt.Println("Property found, enabling property")
			org.Properties.Property[i].Value = value
			found = true
			break
		}
	}

	if !found {
		//set the property
		newProp := types.OrgProperty{}
		newProp.Name = name
		newProp.Value = value

		org.Properties.Property = append(org.Properties.Property, newProp)
	}

	newOrgBody, err := json.Marshal(org)

	u, _ = url.Parse(shared.BaseURL)
	u.Path = path.Join(u.Path, shared.RootArgs.Org)
	_, err = shared.HttpClient(true, u.String(), string(newOrgBody), "PUT")

	return err
}

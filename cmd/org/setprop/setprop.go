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
	_, err = shared.HttpClient(true, u.String(), string(newOrgBody))

	return err
}

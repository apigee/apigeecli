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

package apiclient

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/srinandan/apigeecli/clilog"
)

const apigeecliFile = ".apigeecli"

var usr *user.User

type apigeeCLI struct {
	Token     string `json:"token,omitempty"`
	LastCheck string `json:"lastCheck,omitempty"`
	Org       string `json:"defaultOrg,omitempty"`
	Staging   bool   `json:"staging,omitempty"`
	ProxyUrl  string `json:"proxyUrl,omitempty"`
	Nocheck   bool   `json:"nocheck,omitempty" default:"false"`
}

var cliPref *apigeeCLI //= apigeeCLI{}

func ReadPreferencesFile() (err error) {

	cliPref = new(apigeeCLI)

	usr, err = user.Current()
	if err != nil {
		clilog.Info.Println(err)
		return err
	}

	prefFile, err := ioutil.ReadFile(path.Join(usr.HomeDir, apigeecliFile))
	if err != nil {
		clilog.Info.Println("Cached preferences was not found")
		return err
	}

	err = json.Unmarshal(prefFile, &cliPref)
	clilog.Info.Printf("Token %s, lastCheck: %s", cliPref.Token, cliPref.LastCheck)
	clilog.Info.Printf("DefaultOrg %s", cliPref.Org)
	if err != nil {
		clilog.Info.Printf("Error marshalling: %v\n", err)
		return DeletePreferencesFile()
	}

	if cliPref.Staging {
		UseStaging()
	}

	if cliPref.ProxyUrl != "" {
		SetProxyURL(cliPref.ProxyUrl)
	}

	if cliPref.Org != "" {
		return SetApigeeOrg(cliPref.Org)
	}
	return nil
}

func DeletePreferencesFile() (err error) {
	usr, err = user.Current()
	if err != nil {
		clilog.Info.Println(err)
		return err
	}
	if _, err := os.Stat(path.Join(usr.HomeDir, apigeecliFile)); os.IsNotExist(err) {
		clilog.Info.Println(err)
		return err
	}
	return os.Remove(path.Join(usr.HomeDir, apigeecliFile))
}

func WriteToken(token string) (err error) {
	if IsSkipCache() {
		return nil
	}

	clilog.Info.Println("Cache access token: ", token)
	cliPref.Token = token

	data, err := json.Marshal(&cliPref)
	if err != nil {
		clilog.Info.Printf("Error marshalling: %v\n", err)
		return err
	}
	clilog.Info.Println("Writing ", string(data))
	return WriteByteArrayToFile(path.Join(usr.HomeDir, apigeecliFile), false, data)
}

func GetToken() (token string) {
	return cliPref.Token
}

func GetLastCheck() (lastCheck string) {
	return cliPref.LastCheck
}

func GetNoCheck() bool {
	return cliPref.Nocheck
}

func SetNoCheck(nocheck bool) (err error) {
	clilog.Info.Println("Nocheck set to: ", nocheck)
	cliPref.Nocheck = nocheck

	data, err := json.Marshal(&cliPref)
	if err != nil {
		clilog.Info.Printf("Error marshalling: %v\n", err)
		return err
	}
	clilog.Info.Println("Writing ", string(data))
	return WriteByteArrayToFile(path.Join(usr.HomeDir, apigeecliFile), false, data)
}

func TestAndUpdateLastCheck() (updated bool, err error) {
	currentTime := time.Now()
	currentDate := currentTime.Format("01-02-2006")
	if currentDate == cliPref.LastCheck {
		return true, nil
	}

	cliPref.LastCheck = currentDate

	data, err := json.Marshal(&cliPref)
	if err != nil {
		clilog.Warning.Printf("Error marshalling: %v\n", err)
		return false, err
	}
	clilog.Info.Println("Writing ", string(data))
	if err = WriteByteArrayToFile(path.Join(usr.HomeDir, apigeecliFile), false, data); err != nil {
		return false, err
	}

	return false, nil
}

func GetDefaultOrg() (org string) {
	return cliPref.Org
}

func WriteDefaultOrg(org string) (err error) {
	clilog.Info.Println("Default org: ", org)
	cliPref.Org = org
	data, err := json.Marshal(&cliPref)
	if err != nil {
		clilog.Info.Printf("Error marshalling: %v\n", err)
		return err
	}
	clilog.Info.Println("Writing ", string(data))
	return WriteByteArrayToFile(path.Join(usr.HomeDir, apigeecliFile), false, data)
}

func SetStaging(usestage bool) (err error) {
	if usestage == cliPref.Staging {
		return nil
	}
	cliPref.Staging = usestage
	data, err := json.Marshal(&cliPref)
	if err != nil {
		clilog.Info.Printf("Error marshalling: %v\n", err)
		return err
	}
	clilog.Info.Println("Writing ", string(data))
	return WriteByteArrayToFile(path.Join(usr.HomeDir, apigeecliFile), false, data)
}

func GetStaging() bool {
	return cliPref.Staging
}

func SetProxy(url string) (err error) {
	if url == "" {
		return nil
	}

	cliPref.ProxyUrl = url
	data, err := json.Marshal(&cliPref)
	if err != nil {
		clilog.Info.Printf("Error marshalling: %v\n", err)
		return err
	}
	clilog.Info.Println("Writing ", string(data))
	return WriteByteArrayToFile(path.Join(usr.HomeDir, apigeecliFile), false, data)
}

func GetPreferences() (err error) {
	output, err := json.Marshal(&cliPref)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	PrettyPrint(output)
	return nil
}

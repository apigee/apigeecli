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
}

var cliPref = apigeeCLI{}

func ReadPreferencesFile() (err error) {
	usr, err = user.Current()
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	prefFile, err := ioutil.ReadFile(path.Join(usr.HomeDir, apigeecliFile))
	if err != nil {
		clilog.Info.Println("Cached preferences was not found")
		return err
	}

	err = json.Unmarshal(prefFile, &cliPref)
	if err != nil {
		clilog.Error.Printf("Error marshalling: %v\n", err)
		return err
	}

	return nil
}

func WriteToken(token string) (err error) {
	if IsSkipCache() {
		return nil
	}

	clilog.Info.Println("Cache access token: ", token)
	cliPref.Token = token

	data, err := json.Marshal(&cliPref)
	if err != nil {
		clilog.Error.Printf("Error marshalling: %v\n", err)
		return err
	}
	return WriteByteArrayToFile(path.Join(usr.HomeDir, apigeecliFile), false, data)
}

func GetToken() (token string) {
	return cliPref.Token
}

func GetLastCheck() (lastCheck string) {
	return cliPref.LastCheck
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
		clilog.Error.Printf("Error marshalling: %v\n", err)
		return false, err
	}

	if err = WriteByteArrayToFile(path.Join(usr.HomeDir, apigeecliFile), false, data); err != nil {
		return false, err
	}

	return false, nil
}

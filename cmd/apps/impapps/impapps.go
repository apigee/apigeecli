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

package impapps

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	types "github.com/srinandan/apigeecli/cmd/types"
)

type application struct {
	Name        string            `json:"name,omitempty"`
	Status      string            `json:"status,omitempty"`
	Credentials *[]credential     `json:"credentials,omitempty"`
	DeveloperID *string           `json:"developerId,omitempty"`
	DisplayName string            `json:"displayName,omitempty"`
	Attributes  []types.Attribute `json:"attributes,omitempty"`
	CallbackURL string            `json:"callbackUrl,omitempty"`
	Scopes      []string          `json:"scopes,omitempty"`
}

type credential struct {
	APIProducts    []apiProduct `json:"apiProducts,omitempty"`
	ConsumerKey    string       `json:"consumerKey,omitempty"`
	ConsumerSecret string       `json:"consumerSecret,omitempty"`
	ExpiresAt      int          `json:"expiresAt,omitempty"`
	Status         string       `json:"status,omitempty"`
	Scopes         []string     `json:"scopes,omitempty"`
}

type apiProduct struct {
	Name string `json:"apiproduct,omitempty"`
}

type importCredential struct {
	APIProducts    []string `json:"apiProducts,omitempty"`
	ConsumerKey    string   `json:"consumerKey,omitempty"`
	ConsumerSecret string   `json:"consumerSecret,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
}

//Cmd to import apps
var Cmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing Developer Apps",
	Long:  "Import a file containing Developer Apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		return createApps()
	},
}

var conn int
var file string

func init() {

	Cmd.Flags().StringVarP(&file, "file", "f",
		"", "File containing Developer Apps")
	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = Cmd.MarkFlagRequired("file")
}

func createAsyncApp(app application, wg *sync.WaitGroup) {
	defer wg.Done()

	//importing an app will be a two step process.
	//1. create the app without the credential
	//2. create/import the credential
	u, _ := url.Parse(shared.BaseURL)
	//store the developer and the credential
	developerID := *app.DeveloperID
	credentials := *app.Credentials

	//remove the developer id and credentials from the payload
	app.DeveloperID = nil
	app.Credentials = nil

	out, err := json.Marshal(app)
	if err != nil {
		shared.Error.Fatalln(err)
		return
	}

	u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", developerID, "apps")
	_, err = shared.HttpClient(true, u.String(), string(out))
	if err != nil {
		shared.Error.Fatalln(err)
		return
	}
	u, _ = url.Parse(shared.BaseURL)
	u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", developerID, "apps", app.Name, "keys", "create")
	for _, credential := range credentials {
		//construct a []string for products
		var products []string
		for _, apiProduct := range credential.APIProducts {
			products = append(products, apiProduct.Name)
		}
		//create a new credential
		importCred := importCredential{}
		importCred.APIProducts = products
		importCred.ConsumerKey = credential.ConsumerKey
		importCred.ConsumerSecret = credential.ConsumerSecret
		importCred.Scopes = credential.Scopes

		impCredJSON, err := json.Marshal(importCred)
		if err != nil {
			shared.Error.Fatalln(err)
			return
		}
		_, err = shared.HttpClient(true, u.String(), string(impCredJSON))
		if err != nil {
			shared.Error.Fatalln(err)
			return
		}
		shared.Warning.Println("NOTE: apiProducts are not associated with the app")
	}
	shared.Info.Printf("Completed entity: %s", app.Name)
}

//batch creates a batch of apps to create
func batch(entities []application, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncApp(entity, &bwg)
	}
	bwg.Wait()
}

func createApps() error {

	var pwg sync.WaitGroup

	entities, err := readAppsFile()
	if err != nil {
		shared.Error.Fatalln("Error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	shared.Info.Printf("Found %d apps in the file\n", numEntities)
	shared.Info.Printf("Create apps with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		shared.Info.Printf("Creating batch %d of apps\n", (i + 1))
		go batch(entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		shared.Info.Printf("Creating remaining %d apps\n", remaining)
		go batch(entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

func readAppsFile() ([]application, error) {

	apps := []application{}

	jsonFile, err := os.Open(file)

	if err != nil {
		return apps, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return apps, err
	}

	err = json.Unmarshal(byteValue, &apps)

	if err != nil {
		return apps, err
	}

	return apps, nil
}

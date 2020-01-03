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

package impdev

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

type developer struct {
	EMail      string            `json:"email,omitempty"`
	FirstName  string            `json:"firstName,omitempty"`
	LastName   string            `json:"lastName,omitempty"`
	Attributes []types.Attribute `json:"attributes,omitempty"`
	Username   string            `json:"userName,omitempty"`
}

//Cmd to import developer
var Cmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing App Developers",
	Long:  "Import a file containing App Developers",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers")
		return createDevelopers(u.String())
	},
}

var conn int
var file string

func init() {

	Cmd.Flags().StringVarP(&file, "file", "f",
		"", "File containing App Developers")
	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = Cmd.MarkFlagRequired("file")
}

func createAsyncDeveloper(url string, dev developer, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := json.Marshal(dev)
	if err != nil {
		shared.Error.Fatalln(err)
		return
	}
	_, err = shared.HttpClient(true, url, string(out))
	if err != nil {
		shared.Error.Fatalln(err)
		return
	}

	shared.Info.Printf("Completed entity: %s", dev.EMail)
}

//batch creates a batch of developers to create
func batch(url string, entities []developer, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncDeveloper(url, entity, &bwg)
	}
	bwg.Wait()
}

func createDevelopers(url string) error {

	var pwg sync.WaitGroup

	entities, err := readDevelopersFile()
	if err != nil {
		shared.Error.Fatalln("Error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	shared.Info.Printf("Found %d developers in the file\n", numEntities)
	shared.Info.Printf("Create developers with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		shared.Info.Printf("Creating batch %d of developers\n", (i + 1))
		go batch(url, entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		shared.Info.Printf("Creating remaining %d developers\n", remaining)
		go batch(url, entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

func readDevelopersFile() ([]developer, error) {

	devs := []developer{}

	jsonFile, err := os.Open(file)

	if err != nil {
		return devs, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return devs, err
	}

	err = json.Unmarshal(byteValue, &devs)

	if err != nil {
		return devs, err
	}

	return devs, nil
}

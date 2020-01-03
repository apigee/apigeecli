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

package impprod

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

type product struct {
	Name         string            `json:"name,omitempty"`
	DisplayName  string            `json:"displayName,omitempty"`
	ApprovalType string            `json:"approvalType,omitempty"`
	Attributes   []types.Attribute `json:"attributes,omitempty"`
	APIResources []string          `json:"apiResources,omitempty"`
	Environments []string          `json:"environments,omitempty"`
	Proxies      []string          `json:"proxies,omitempty"`
}

//Cmd to import products
var Cmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing API products",
	Long:  "Import a file containing API products",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apiproducts")
		return createProducts(u.String())
	},
}

var conn int
var file string

func init() {

	Cmd.Flags().StringVarP(&file, "file", "f",
		"", "File containing API Products")
	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = Cmd.MarkFlagRequired("file")
}

func createAsyncProduct(url string, entity product, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := json.Marshal(entity)
	if err != nil {
		shared.Error.Fatalln(err)
		return
	}
	_, err = shared.HttpClient(true, url, string(out))
	if err != nil {
		shared.Error.Fatalln(err)
		return
	}
	shared.Info.Printf("Completed entity: %s", entity.Name)
}

//batch creates a batch of products to create
func batch(url string, entities []product, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncProduct(url, entity, &bwg)
	}
	bwg.Wait()
}

func createProducts(url string) error {

	var pwg sync.WaitGroup

	entities, err := readProductsFile()
	if err != nil {
		shared.Error.Fatalln("Error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	shared.Info.Printf("Found %d products in the file\n", numEntities)
	shared.Info.Printf("Create products with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		shared.Info.Printf("Creating batch %d of products\n", (i + 1))
		go batch(url, entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		shared.Info.Printf("Creating remaining %d products\n", remaining)
		go batch(url, entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

func readProductsFile() ([]product, error) {

	products := []product{}

	jsonFile, err := os.Open(file)

	if err != nil {
		return products, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return products, err
	}

	err = json.Unmarshal(byteValue, &products)

	if err != nil {
		return products, err
	}

	return products, nil
}

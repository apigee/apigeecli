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

package products

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/clilog"
)

type apiProducts struct {
	APIProducts []apiProduct `json:"apiProduct,omitempty"`
}

type apiProduct struct {
	Name string `json:"name,omitempty"`
}

type product struct {
	Name          string      `json:"name,omitempty"`
	DisplayName   string      `json:"displayName,omitempty"`
	Description   string      `json:"description,omitempty"`
	ApprovalType  string      `json:"approvalType,omitempty"`
	Attributes    []attribute `json:"attributes,omitempty"`
	APIResources  []string    `json:"apiResources,omitempty"`
	Environments  []string    `json:"environments,omitempty"`
	Proxies       []string    `json:"proxies,omitempty"`
	Quota         string      `json:"quota,omitempty"`
	QuotaInterval string      `json:"quotaInterval,omitempty"`
	QuotaTimeUnit string      `json:"quotaTimeUnit,omitempty"`
	Scopes        []string    `json:"scopes,omitempty"`
}

//attribute to used to hold custom attributes for entities
type attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

//Create
func Create(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	product := []string{}

	product = append(product, "\"name\":\""+name+"\"")

	if displayName == "" {
		product = append(product, "\"displayName\":\""+name+"\"")
	} else {
		product = append(product, "\"displayName\":\""+displayName+"\"")
	}

	if description != "" {
		product = append(product, "\"description\":\""+description+"\"")
	}
	product = append(product, "\"environments\":[\""+getArrayStr(environments)+"\"]")
	product = append(product, "\"proxies\":[\""+getArrayStr(proxies)+"\"]")

	if len(scopes) > 0 {
		product = append(product, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")
	}

	product = append(product, "\"approvalType\":\""+approval+"\"")

	if quota != "" {
		product = append(product, "\"quota\":\""+quota+"\"")
	}
	if quotaInterval != "" {
		product = append(product, "\"quotaInterval\":\""+quotaInterval+"\"")
	}
	if quotaUnit != "" {
		product = append(product, "\"quotaTimeUnit\":\""+quotaUnit+"\"")
	}
	if len(attrs) != 0 {
		attributes := []string{}
		for key, value := range attrs {
			attributes = append(attributes, "{\"name\":\""+key+"\",\"value\":\""+value+"\"}")
		}
		attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
		product = append(product, attributesStr)
	}

	payload := "{" + strings.Join(product, ",") + "}"
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//Update
func Update(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	apiclient.SetPrintOutput(false)
	respBody, err = Get(name)
	if err != nil {
		return nil, err
	}
	apiclient.SetPrintOutput(true)

	p := product{}

	err = json.Unmarshal(respBody, &p)
	if err != nil {
		return nil, err
	}

	if displayName != "" {
		p.DisplayName = displayName
	}

	if description != "" {
		p.Description = description
	}

	if quota != "" {
		p.Quota = quota
	}

	if quotaInterval != "" {
		p.QuotaInterval = quotaInterval
	}

	if quotaUnit != "" {
		p.QuotaTimeUnit = quotaUnit
	}

	if len(environments) > 0 {
		p.Environments = append(p.Environments, environments...)
	}

	if len(proxies) > 0 {
		p.Proxies = append(p.Proxies, proxies...)
	}

	if len(scopes) > 0 {
		p.Scopes = append(p.Scopes, scopes...)
	}

	if len(attrs) > 0 {
		//create new attributes
		for k, v := range attrs {
			a := attribute{}
			a.Name = k
			a.Value = v
			p.Attributes = append(p.Attributes, a)
		}
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(payload), "PUT")

	return respBody, err
}

//UpdateAttribute
func UpdateAttribute(name string, key string, value string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name, "attributes", key)
	payload := "{ \"name\":\"" + key + "\",\"value\":\"" + value + "\"}"
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//DeleteAttribute
func DeleteAttribute(name string, key string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name, "attributes", key)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//GetAttribute
func GetAttribute(name string, key string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name, "attributes", key)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//ListAttributes
func ListAttributes(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name, "attributes")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//List
func List(count int, expand bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts")
	q := u.Query()
	if expand {
		q.Set("expand", "true")
	} else {
		q.Set("expand", "false")
	}
	if count != -1 {
		q.Set("count", strconv.Itoa(count))
	}
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//Export
func Export(conn int) (payload [][]byte, err error) {
	//parent workgroup
	var pwg sync.WaitGroup
	var mu sync.Mutex
	const entityType = "apiproducts"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), entityType)
	//don't print to sysout
	respBody, err := apiclient.HttpClient(false, u.String())
	if err != nil {
		return apiclient.GetEntityPayloadList(), nil
	}

	var products = apiProducts{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		return apiclient.GetEntityPayloadList(), nil
	}

	numProd := len(products.APIProducts)
	clilog.Info.Printf("Found %d products in the org\n", numProd)
	clilog.Info.Printf("Exporting products with %d connections\n", conn)

	numOfLoops, remaining := numProd/conn, numProd%conn

	//ensure connections aren't greater than products
	if conn > numProd {
		conn = numProd
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Exporting batch %d of products\n", (i + 1))
		go batchExport(products.APIProducts[start:end], entityType, &pwg, &mu)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Exporting remaining %d products\n", remaining)
		go batchExport(products.APIProducts[start:numProd], entityType, &pwg, &mu)
		pwg.Wait()
	}

	return apiclient.GetEntityPayloadList(), nil
}

//Import
func Import(conn int, filePath string) (err error) {
	var pwg sync.WaitGroup
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts")

	entities, err := readProductsFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	clilog.Info.Printf("Found %d products in the file\n", numEntities)
	clilog.Info.Printf("Create products with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Info.Printf("Creating batch %d of products\n", (i + 1))
		go batchImport(u.String(), entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Info.Printf("Creating remaining %d products\n", remaining)
		go batchImport(u.String(), entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

//batch created a batch of products to query
func batchExport(entities []apiProduct, entityType string, pwg *sync.WaitGroup, mu *sync.Mutex) {
	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		u, _ := url.Parse(apiclient.BaseURL)
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), entityType, url.PathEscape(entity.Name))
		go apiclient.GetAsyncEntity(u.String(), &bwg, mu)
	}
	bwg.Wait()
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}

//batch creates a batch of products to create
func batchImport(url string, entities []product, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncProduct(url, entity, &bwg)
	}
	bwg.Wait()
}

func createAsyncProduct(url string, entity product, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := json.Marshal(entity)
	if err != nil {
		clilog.Error.Println(err)
		return
	}
	_, err = apiclient.HttpClient(apiclient.GetPrintOutput(), url, string(out))
	if err != nil {
		clilog.Error.Println(err)
		return
	}
	clilog.Info.Printf("Completed entity: %s", entity.Name)
}

func readProductsFile(filePath string) ([]product, error) {

	products := []product{}

	jsonFile, err := os.Open(filePath)

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

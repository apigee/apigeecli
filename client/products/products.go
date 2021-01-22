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
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
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

const proxyOperationConfigType = "proxy"
const remoteServiceOperationConfigType = "remoteservice"

type product struct {
	Name           string          `json:"name,omitempty"`
	DisplayName    string          `json:"displayName,omitempty"`
	Description    string          `json:"description,omitempty"`
	ApprovalType   string          `json:"approvalType,omitempty"`
	Attributes     []attribute     `json:"attributes,omitempty"`
	APIResources   []string        `json:"apiResources,omitempty"`
	OperationGroup *operationGroup `json:"operationGroup,omitempty"`
	Environments   []string        `json:"environments,omitempty"`
	Proxies        []string        `json:"proxies,omitempty"`
	Quota          string          `json:"quota,omitempty"`
	QuotaInterval  string          `json:"quotaInterval,omitempty"`
	QuotaTimeUnit  string          `json:"quotaTimeUnit,omitempty"`
	Scopes         []string        `json:"scopes,omitempty"`
}

type operationGroup struct {
	OperationConfigs    []operationConfig `json:"operationConfigs,omitempty"`
	OperationConfigType string            `json:"operationConfigType,omitempty"`
}

type operationConfig struct {
	APISource  string      `json:"apiSource,omitempty"`
	Operations []operation `json:"operations,omitempty"`
	Quota      *quota      `json:"quota,omitempty"`
	Attributes []attribute `json:"attributes,omitempty"`
}

type operation struct {
	Resource string   `json:"resource,omitempty"`
	Methods  []string `json:"methods,omitempty"`
}

type quota struct {
	Limit    string `json:"limit,omitempty"`
	Interval string `json:"interval,omitempty"`
	TimeUnit string `json:"timeUnit,omitempty"`
}

//attribute to used to hold custom attributes for entities
type attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func CreateLegacy(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	return createProduct(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, proxies, scopes, nil, true, attrs, true)
}

func Create(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	return createProduct(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, proxies, scopes, nil, true, attrs, false)
}

func CreateRemoteServiceOperationGroup(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	return createProduct(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, proxies, scopes, nil, false, attrs, false)
}

func CreateProxyOperationGroup(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, scopes []string, operationGrp []byte, attrs map[string]string) (respBody []byte, err error) {
	return createProduct(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, nil, scopes, operationGrp, true, attrs, false)
}

func UpdateLegacy(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	return updateProduct(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, proxies, scopes, nil, true, attrs, true)
}

func Update(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	return updateProduct(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, proxies, scopes, nil, true, attrs, false)
}

func UpdateRemoteServiceOperationGroup(name string, services []string) (respBody []byte, err error) {
	return updateProduct(name, "", "", "", "", "", "", nil, services, nil, nil, false, nil, false)
}

func UpdateProxyOperationGroup(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, scopes []string, operationGrp []byte, attrs map[string]string) (respBody []byte, err error) {
	return updateProduct(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, nil, scopes, operationGrp, true, attrs, false)
}

//createProduct
func createProduct(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, operationGrp []byte, proxyOperationGroup bool, attrs map[string]string, legacy bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	var Product = new(product)
	var OperationGroup = new(operationGroup)

	if len(proxies) > 0 && len(operationGrp) > 0 {
		return nil, fmt.Errorf("A product cannot have proxies and operation group")
	}

	Product.Name = name
	Product.ApprovalType = approval

	if len(operationGrp) > 0 {
		err = json.Unmarshal(operationGrp, OperationGroup)
		if err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		if reflect.DeepEqual(*OperationGroup, operationGroup{}) {
			return nil, fmt.Errorf("can't unmarshal json to OperationGroup")
		}
		Product.OperationGroup = OperationGroup
	}

	if displayName == "" {
		Product.DisplayName = name
	} else {
		Product.DisplayName = displayName
	}

	if description != "" {
		Product.Description = description
	}

	if len(environments) > 0 {
		Product.Environments = environments
	}

	if len(proxies) > 0 {
		if legacy {
			Product.Proxies = proxies
		} else {
			OperationConfigs := []operationConfig{}
			Operations := []operation{}
			Operation := operation{}
			Operations = append(Operations, Operation)

			for _, proxy := range proxies {
				OperationConfig := operationConfig{}
				OperationConfig.APISource = proxy
				OperationConfig.Operations = Operations
				OperationConfigs = append(OperationConfigs, OperationConfig)
			}
			OperationGroup.OperationConfigs = OperationConfigs
			if proxyOperationGroup {
				OperationGroup.OperationConfigType = proxyOperationConfigType
			} else {
				OperationGroup.OperationConfigType = remoteServiceOperationConfigType
			}
			Product.OperationGroup = OperationGroup
		}
	}

	if len(scopes) > 0 {
		Product.Scopes = scopes
	}

	if quota != "" {
		Product.Quota = quota
	}
	if quotaInterval != "" {
		Product.QuotaInterval = quotaInterval
	}

	if quotaUnit != "" {
		Product.QuotaTimeUnit = quotaUnit
	}

	if len(attrs) > 0 {
		//create new attributes
		for k, v := range attrs {
			a := attribute{}
			a.Name = k
			a.Value = v
			Product.Attributes = append(Product.Attributes, a)
		}
	}

	payload, err := json.Marshal(Product)
	if err != nil {
		clilog.Info.Println(err)
		return nil, err
	}

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(payload))
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

//updateProduct
func updateProduct(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, operationGrp []byte, proxyOperationGroup bool, attrs map[string]string, legacy bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	if len(proxies) > 0 && len(operationGrp) > 0 {
		return nil, fmt.Errorf("A product cannot have proxies and operation group")
	}

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
		if legacy {
			p.Proxies = append(p.Proxies, proxies...)
		} else {
			var OperationGroup = new(operationGroup)
			OperationConfigs := []operationConfig{}
			Operations := []operation{}
			Operation := operation{}
			Operations = append(Operations, Operation)

			for _, proxy := range proxies {
				OperationConfig := operationConfig{}
				OperationConfig.APISource = proxy
				OperationConfig.Operations = Operations
				OperationConfigs = append(OperationConfigs, OperationConfig)
			}
			OperationGroup.OperationConfigs = OperationConfigs
			if proxyOperationGroup {
				OperationGroup.OperationConfigType = proxyOperationConfigType
			} else {
				OperationGroup.OperationConfigType = remoteServiceOperationConfigType
			}
			p.OperationGroup = OperationGroup
		}
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
		clilog.Info.Println(err)
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
		return apiclient.GetEntityPayloadList(), err
	}

	var products = apiProducts{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
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

	payload = make([][]byte, len(apiclient.GetEntityPayloadList()))
	copy(payload, apiclient.GetEntityPayloadList())
	apiclient.ClearEntityPayloadList()
	return payload, nil
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

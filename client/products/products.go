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
	Name                  string                 `json:"name,omitempty"`
	DisplayName           string                 `json:"displayName,omitempty"`
	Description           string                 `json:"description,omitempty"`
	ApprovalType          string                 `json:"approvalType,omitempty"`
	Attributes            []attribute            `json:"attributes,omitempty"`
	APIResources          []string               `json:"apiResources,omitempty"`
	OperationGroup        *operationGroup        `json:"operationGroup,omitempty"`
	GraphQLOperationGroup *graphqlOperationGroup `json:"graphqlOperationGroup,omitempty"`
	Environments          []string               `json:"environments,omitempty"`
	Proxies               []string               `json:"proxies,omitempty"`
	Quota                 string                 `json:"quota,omitempty"`
	QuotaInterval         string                 `json:"quotaInterval,omitempty"`
	QuotaTimeUnit         string                 `json:"quotaTimeUnit,omitempty"`
	Scopes                []string               `json:"scopes,omitempty"`
}

type operationGroup struct {
	OperationConfigs    []operationConfig `json:"operationConfigs,omitempty"`
	OperationConfigType string            `json:"operationConfigType,omitempty"`
}

type graphqlOperationGroup struct {
	OperationConfigs    []graphQLOperationConfig `json:"operationConfigs,omitempty"`
	OperationConfigType string                   `json:"operationConfigType,omitempty"`
}

type operationConfig struct {
	APISource  string      `json:"apiSource,omitempty"`
	Operations []operation `json:"operations,omitempty"`
	Quota      *quota      `json:"quota,omitempty"`
	Attributes []attribute `json:"attributes,omitempty"`
}

type graphQLOperationConfig struct {
	APISource  string             `json:"apiSource,omitempty"`
	Operations []graphQLoperation `json:"operations,omitempty"`
	Quota      *quota             `json:"quota,omitempty"`
	Attributes []attribute        `json:"attributes,omitempty"`
}

type operation struct {
	Resource string   `json:"resource,omitempty"`
	Methods  []string `json:"methods,omitempty"`
}

type graphQLoperation struct {
	OperationTypes []string `json:"operationTypes,omitempty"`
	Operation      string   `json:"operation,omitempty"`
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

type ProductSettings struct {
	Name            string
	Description     string
	Approval        string
	DisplayName     string
	Quota           string
	QuotaInterval   string
	QuotaUnit       string
	Environments    []string
	Proxies         []string
	Scopes          []string
	Attrs           map[string]string
	OperationGrp    []byte
	GqlOperationGrp []byte
	Legacy          bool
}

func CreateLegacy(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	productSettings := ProductSettings{}
	productSettings.Name = name
	productSettings.Description = description
	productSettings.Approval = approval
	productSettings.DisplayName = displayName
	productSettings.Quota = quota
	productSettings.QuotaInterval = quotaInterval
	productSettings.QuotaUnit = quotaUnit
	productSettings.Environments = environments
	productSettings.Proxies = proxies
	productSettings.Scopes = scopes
	productSettings.Attrs = attrs
	productSettings.Legacy = true

	return createProduct(productSettings)
}

func CreateProxyOperationGroup(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, scopes []string, operationGrp []byte, gqlOperationGrp []byte, attrs map[string]string) (respBody []byte, err error) {
	productSettings := ProductSettings{}
	productSettings.Name = name
	productSettings.Description = description
	productSettings.Approval = approval
	productSettings.DisplayName = displayName
	productSettings.Quota = quota
	productSettings.QuotaInterval = quotaInterval
	productSettings.QuotaUnit = quotaUnit
	productSettings.Environments = environments
	productSettings.Scopes = scopes
	productSettings.Attrs = attrs

	productSettings.OperationGrp = operationGrp
	productSettings.GqlOperationGrp = gqlOperationGrp

	productSettings.Legacy = false

	return createProduct(productSettings)
}

func UpdateLegacy(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, proxies []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	productSettings := ProductSettings{}
	productSettings.Name = name
	productSettings.Description = description
	productSettings.Approval = approval
	productSettings.DisplayName = displayName
	productSettings.Quota = quota
	productSettings.QuotaInterval = quotaInterval
	productSettings.QuotaUnit = quotaUnit
	productSettings.Environments = environments
	productSettings.Proxies = proxies
	productSettings.Scopes = scopes
	productSettings.Attrs = attrs
	productSettings.Legacy = true

	return updateProduct(productSettings)
}

func UpdateProxyOperationGroup(name string, description string, approval string, displayName string, quota string, quotaInterval string, quotaUnit string, environments []string, scopes []string, operationGrp []byte, gqlOperationGrp []byte, attrs map[string]string) (respBody []byte, err error) {
	productSettings := ProductSettings{}
	productSettings.Name = name
	productSettings.Description = description
	productSettings.Approval = approval
	productSettings.DisplayName = displayName
	productSettings.Quota = quota
	productSettings.QuotaInterval = quotaInterval
	productSettings.QuotaUnit = quotaUnit
	productSettings.Environments = environments
	productSettings.Scopes = scopes
	productSettings.Attrs = attrs

	productSettings.OperationGrp = operationGrp
	productSettings.GqlOperationGrp = gqlOperationGrp

	productSettings.Legacy = false

	return updateProduct(productSettings)
}

//createProduct
func createProduct(productSettings ProductSettings) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	var Product = new(product)
	var OperationGroup = new(operationGroup)
	var GraphqlOperationGroup = new(graphqlOperationGroup)

	Product.Name = productSettings.Name
	Product.ApprovalType = productSettings.Approval

	if len(productSettings.OperationGrp) > 0 {
		err = json.Unmarshal(productSettings.OperationGrp, OperationGroup)
		if err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		if reflect.DeepEqual(*OperationGroup, operationGroup{}) {
			return nil, fmt.Errorf("can't unmarshal json to OperationGroup")
		}
		Product.OperationGroup = OperationGroup
	}

	if len(productSettings.GqlOperationGrp) > 0 {
		err = json.Unmarshal(productSettings.GqlOperationGrp, GraphqlOperationGroup)
		if err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		if reflect.DeepEqual(*GraphqlOperationGroup, graphqlOperationGroup{}) {
			return nil, fmt.Errorf("can't unmarshal json to GraphqlOperationGroup")
		}
		Product.GraphQLOperationGroup = GraphqlOperationGroup
	}

	if productSettings.DisplayName == "" {
		Product.DisplayName = productSettings.Name
	} else {
		Product.DisplayName = productSettings.DisplayName
	}

	if productSettings.Description != "" {
		Product.Description = productSettings.Description
	}

	if len(productSettings.Environments) > 0 {
		Product.Environments = productSettings.Environments
	}

	if len(productSettings.Proxies) > 0 {
		if productSettings.Legacy {
			Product.Proxies = productSettings.Proxies
		} else if len(productSettings.OperationGrp) == 0 {
			OperationConfigs := []operationConfig{}
			Operations := []operation{}
			Operation := operation{}
			Operations = append(Operations, Operation)

			for _, proxy := range productSettings.Proxies {
				OperationConfig := operationConfig{}
				OperationConfig.APISource = proxy
				OperationConfig.Operations = Operations
				OperationConfigs = append(OperationConfigs, OperationConfig)
			}
			OperationGroup.OperationConfigs = OperationConfigs
			OperationGroup.OperationConfigType = proxyOperationConfigType
			Product.OperationGroup = OperationGroup
		}
	}

	if len(productSettings.Scopes) > 0 {
		Product.Scopes = productSettings.Scopes
	}

	if productSettings.Quota != "" {
		Product.Quota = productSettings.Quota
	}
	if productSettings.QuotaInterval != "" {
		Product.QuotaInterval = productSettings.QuotaInterval
	}

	if productSettings.QuotaUnit != "" {
		Product.QuotaTimeUnit = productSettings.QuotaUnit
	}

	if len(productSettings.Attrs) > 0 {
		//create new attributes
		for k, v := range productSettings.Attrs {
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
func updateProduct(productSettings ProductSettings) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	var OperationGroup = new(operationGroup)
	var GraphqlOperationGroup = new(graphqlOperationGroup)

	apiclient.SetPrintOutput(false)
	respBody, err = Get(productSettings.Name)
	if err != nil {
		return nil, err
	}
	apiclient.SetPrintOutput(true)

	p := product{}

	err = json.Unmarshal(respBody, &p)
	if err != nil {
		return nil, err
	}

	if productSettings.DisplayName != "" {
		p.DisplayName = productSettings.DisplayName
	}

	if productSettings.Description != "" {
		p.Description = productSettings.Description
	}

	if productSettings.Quota != "" {
		p.Quota = productSettings.Quota
	}

	if productSettings.QuotaInterval != "" {
		p.QuotaInterval = productSettings.QuotaInterval
	}

	if productSettings.QuotaUnit != "" {
		p.QuotaTimeUnit = productSettings.QuotaUnit
	}

	if len(productSettings.Environments) > 0 {
		p.Environments = append(p.Environments, productSettings.Environments...)
	}

	if len(productSettings.OperationGrp) > 0 {
		err = json.Unmarshal(productSettings.OperationGrp, OperationGroup)
		if err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		if reflect.DeepEqual(*OperationGroup, operationGroup{}) {
			return nil, fmt.Errorf("can't unmarshal json to OperationGroup")
		}
		//check to see the operation config type is the same
		if OperationGroup.OperationConfigType != p.OperationGroup.OperationConfigType {
			return nil, fmt.Errorf("updated operationConfigType must match the existing operationConfigType - ", OperationGroup.OperationConfigType)
		}
		p.OperationGroup.OperationConfigs = append(p.OperationGroup.OperationConfigs, OperationGroup.OperationConfigs...)
	}

	if len(productSettings.GqlOperationGrp) > 0 {
		err = json.Unmarshal(productSettings.GqlOperationGrp, GraphqlOperationGroup)
		if err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		if reflect.DeepEqual(*GraphqlOperationGroup, graphqlOperationGroup{}) {
			return nil, fmt.Errorf("can't unmarshal json to GraphqlOperationGroup")
		}
		//check to see the operation config type is the same
		if GraphqlOperationGroup.OperationConfigType != p.GraphQLOperationGroup.OperationConfigType {
			return nil, fmt.Errorf("updated operationConfigType must match the existing operationConfigType - ", GraphqlOperationGroup.OperationConfigType)
		}
		p.GraphQLOperationGroup.OperationConfigs = append(p.GraphQLOperationGroup.OperationConfigs, GraphqlOperationGroup.OperationConfigs...)
	}

	if len(productSettings.Proxies) > 0 {
		if productSettings.Legacy {
			p.Proxies = append(p.Proxies, productSettings.Proxies...)
		} else if len(productSettings.OperationGrp) == 0 {
			var OperationGroup = new(operationGroup)
			OperationConfigs := []operationConfig{}
			Operations := []operation{}
			Operation := operation{}
			Operations = append(Operations, Operation)

			for _, proxy := range productSettings.Proxies {
				OperationConfig := operationConfig{}
				OperationConfig.APISource = proxy
				OperationConfig.Operations = Operations
				OperationConfigs = append(OperationConfigs, OperationConfig)
			}
			OperationGroup.OperationConfigs = OperationConfigs
			OperationGroup.OperationConfigType = proxyOperationConfigType
			p.OperationGroup = OperationGroup
		}
	}

	if len(productSettings.Scopes) > 0 {
		p.Scopes = append(p.Scopes, productSettings.Scopes...)
	}

	if len(productSettings.Attrs) > 0 {
		//create new attributes
		for k, v := range productSettings.Attrs {
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

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", productSettings.Name)
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

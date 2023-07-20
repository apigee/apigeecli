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
	"errors"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"internal/apiclient"

	"internal/clilog"
)

type apiProducts struct {
	APIProduct []APIProduct `json:"apiProduct,omitempty"`
}

type Action uint8

const (
	CREATE Action = iota
	UPDATE
	UPSERT
)

type APIProduct struct {
	Name                  string                 `json:"name,omitempty"`
	DisplayName           string                 `json:"displayName,omitempty"`
	Description           string                 `json:"description,omitempty"`
	ApprovalType          string                 `json:"approvalType,omitempty"`
	Attributes            []Attribute            `json:"attributes,omitempty"`
	APIResources          []string               `json:"apiResources,omitempty"`
	OperationGroup        *OperationGroup        `json:"operationGroup,omitempty"`
	GraphQLOperationGroup *GraphqlOperationGroup `json:"graphqlOperationGroup,omitempty"`
	GrpcOperationGroup    *GrpcOperationGroup    `json:"grpcOperationGroup,omitempty"`
	Environments          []string               `json:"environments,omitempty"`
	Proxies               []string               `json:"proxies,omitempty"`
	Quota                 string                 `json:"quota,omitempty"`
	QuotaInterval         string                 `json:"quotaInterval,omitempty"`
	QuotaTimeUnit         string                 `json:"quotaTimeUnit,omitempty"`
	Scopes                []string               `json:"scopes,omitempty"`
}

type OperationGroup struct {
	OperationConfigs    []operationConfig `json:"operationConfigs,omitempty"`
	OperationConfigType string            `json:"operationConfigType,omitempty"`
}

type GraphqlOperationGroup struct {
	OperationConfigs    []graphQLOperationConfig `json:"operationConfigs,omitempty"`
	OperationConfigType string                   `json:"operationConfigType,omitempty"`
}

type GrpcOperationGroup struct {
	OperationConfigs    []grpcOperationConfig `json:"operationConfigs,omitempty"`
	OperationConfigType string                `json:"operationConfigType,omitempty"`
}

type operationConfig struct {
	APISource  string      `json:"apiSource,omitempty"`
	Operations []operation `json:"operations,omitempty"`
	Quota      *quota      `json:"quota,omitempty"`
	Attributes []Attribute `json:"attributes,omitempty"`
}

type graphQLOperationConfig struct {
	APISource  string             `json:"apiSource,omitempty"`
	Operations []graphQLoperation `json:"operations,omitempty"`
	Quota      *quota             `json:"quota,omitempty"`
	Attributes []Attribute        `json:"attributes,omitempty"`
}

type grpcOperationConfig struct {
	APISource  string      `json:"apiSource,omitempty"`
	Service    string      `json:"service,omitempty"`
	Methods    []string    `json:"methods,omitempty"`
	Quota      *quota      `json:"quota,omitempty"`
	Attributes []Attribute `json:"attributes,omitempty"`
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

// attribute to used to hold custom attributes for entities
type Attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func Create(p APIProduct) (respBody []byte, err error) {
	return upsert(p, CREATE)
}

func Update(p APIProduct) (respBody []byte, err error) {
	return upsert(p, UPDATE)
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// upsert - use Action to control if upsert is enabled
func upsert(p APIProduct, a Action) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	var createNew bool // default false

	switch a {
	case CREATE:
		createNew = true
	case UPDATE:
		createNew = false
	case UPSERT:
		apiclient.ClientPrintHttpResponse.Set(false)
		_, err = Get(p.Name)
		if err != nil {
			createNew = true // product does not exist
		}
		apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	}

	payload, err := json.Marshal(p)
	if err != nil {
		clilog.Debug.Println(err)
		return nil, err
	}

	if createNew {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts")
		respBody, err = apiclient.HttpClient(u.String(), string(payload))
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", p.Name)
		respBody, err = apiclient.HttpClient(u.String(), string(payload), "PUT")
	}

	return respBody, err
}

// UpdateAttribute
func UpdateAttribute(name string, key string, value string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name, "attributes", key)
	payload := "{ \"name\":\"" + key + "\",\"value\":\"" + value + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// DeleteAttribute
func DeleteAttribute(name string, key string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name, "attributes", key)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// GetAttribute
func GetAttribute(name string, key string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name, "attributes", key)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListAttributes
func ListAttributes(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apiproducts", name, "attributes")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// List
func List(count int, startKey string, expand bool) (respBody []byte, err error) {
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
	if startKey != "" {
		q.Set("startKey", startKey)
	}

	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String())

	return respBody, err
}

// ListFilter
func ListFilter(filter map[string]string) (respBody []byte, err error) {
	maxProducts := 1000
	nextPage := true
	startKey := ""
	allprds := apiProducts{}
	outprds := apiProducts{}

	apiclient.ClientPrintHttpResponse.Set(false)

	for nextPage {
		pageResp, err := List(maxProducts, startKey, true)
		if err != nil {
			return nil, err
		}
		prds := apiProducts{}
		err = json.Unmarshal(pageResp, &prds)
		if err != nil {
			return nil, err
		}

		startKey = prds.APIProduct[len(prds.APIProduct)-1].Name

		allprds.APIProduct = append(allprds.APIProduct, prds.APIProduct...)

		// if there is only one item in the list, the there are no more products to fetch
		if len(prds.APIProduct) == 1 {
			nextPage = false
		}
	}

	if filter["proxy"] != "" {
		for _, p := range allprds.APIProduct {
			if p.OperationGroup != nil && len(p.OperationGroup.OperationConfigs) > 0 {
				if p.OperationGroup.OperationConfigType == "proxy" {
					for _, o := range p.OperationGroup.OperationConfigs {
						if o.APISource == filter["proxy"] {
							outprds.APIProduct = append(outprds.APIProduct, p)
						}
					}
				}
			}
			if p.GraphQLOperationGroup != nil && len(p.GraphQLOperationGroup.OperationConfigs) > 0 {
				for _, o := range p.GraphQLOperationGroup.OperationConfigs {
					if o.APISource == filter["proxy"] {
						outprds.APIProduct = append(outprds.APIProduct, p)
					}
				}
			}
			if p.GrpcOperationGroup != nil && len(p.GraphQLOperationGroup.OperationConfigs) > 0 {
				for _, o := range p.GrpcOperationGroup.OperationConfigs {
					if o.APISource == filter["proxy"] {
						outprds.APIProduct = append(outprds.APIProduct, p)
					}
				}
			}
		}
	} else {
		outprds = allprds
	}

	respBody, err = json.Marshal(outprds)
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	_ = apiclient.PrettyPrint("json", respBody)

	return respBody, err
}

// Export
func Export(conn int) (payload [][]byte, err error) {
	// parent workgroup
	var pwg sync.WaitGroup
	var mu sync.Mutex
	const entityType = "apiproducts"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), entityType)
	// don't print to sysout
	apiclient.ClientPrintHttpResponse.Set(false)
	respBody, err := apiclient.HttpClient(u.String())
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	products := apiProducts{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		return apiclient.GetEntityPayloadList(), err
	}

	numProd := len(products.APIProduct)
	clilog.Debug.Printf("Found %d products in the org\n", numProd)
	clilog.Debug.Printf("Exporting products with %d connections\n", conn)

	numOfLoops, remaining := numProd/conn, numProd%conn

	// ensure connections aren't greater than products
	if conn > numProd {
		conn = numProd
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Debug.Printf("Exporting batch %d of products\n", (i + 1))
		go batchExport(products.APIProduct[start:end], entityType, &pwg, &mu)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Debug.Printf("Exporting remaining %d products\n", remaining)
		go batchExport(products.APIProduct[start:numProd], entityType, &pwg, &mu)
		pwg.Wait()
	}

	payload = make([][]byte, len(apiclient.GetEntityPayloadList()))
	copy(payload, apiclient.GetEntityPayloadList())
	apiclient.ClearEntityPayloadList()
	return payload, nil
}

// Import
func Import(conn int, filePath string, upsert bool) (err error) {
	entities, err := readProductsFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	clilog.Debug.Printf("Found %d products in the file\n", numEntities)
	clilog.Debug.Printf("Create products with %d connections\n", conn)

	jobChan := make(chan APIProduct)
	errChan := make(chan error)

	fanOutWg := sync.WaitGroup{}
	fanInWg := sync.WaitGroup{}

	errs := []string{}
	fanInWg.Add(1)
	go func() {
		defer fanInWg.Done()
		for {
			newErr, ok := <-errChan
			if !ok {
				return
			}
			errs = append(errs, newErr.Error())
		}
	}()

	for i := 0; i < conn; i++ {
		fanOutWg.Add(1)
		go importAPIProduct(&fanOutWg, upsert, jobChan, errChan)
	}

	for _, entity := range entities {
		jobChan <- entity
	}
	close(jobChan)
	fanOutWg.Wait()
	close(errChan)
	fanInWg.Wait()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// batch created a batch of products to query
func batchExport(entities []APIProduct, entityType string, pwg *sync.WaitGroup, mu *sync.Mutex) {
	defer pwg.Done()
	// batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		u, _ := url.Parse(apiclient.BaseURL)
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), entityType, url.PathEscape(entity.Name))
		go apiclient.GetAsyncEntity(u.String(), &bwg, mu)
	}
	bwg.Wait()
}

// importAPIProduct
func importAPIProduct(wg *sync.WaitGroup, upsertAction bool, jobs <-chan APIProduct, errs chan<- error) {
	var err error
	defer wg.Done()
	for {
		job, ok := <-jobs
		if !ok {
			return
		}
		if upsertAction {
			_, err = upsert(job, UPSERT)
		} else {
			_, err = upsert(job, CREATE)
		}
		if err != nil {
			errs <- err
			continue
		}
	}
}

func readProductsFile(filePath string) ([]APIProduct, error) {
	products := []APIProduct{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return products, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return products, err
	}

	err = json.Unmarshal(byteValue, &products)

	if err != nil {
		return products, err
	}

	return products, nil
}

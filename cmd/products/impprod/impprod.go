package impprod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	types "github.com/srinandan/apigeecli/cmd/types"
)

type Product struct {
	Name         string            `json:"name,omitempty"`
	DisplayName  string            `json:"displayName,omitempty"`
	ApprovalType string            `json:"approvalType,omitempty"`
	Attributes   []types.Attribute `json:"attributes,omitempty"`
	APIResources []string          `json:"apiResources,omitempty"`
	Environments []string          `json:"environments,omitempty"`
	Proxies      []string          `json:"proxies,omitempty"`
}

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

func createAsyncProduct(url string, product Product, wg *sync.WaitGroup, errChan chan<- *types.ImportError) {
	defer wg.Done()
	out, err := json.Marshal(product)
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}
	_, err = shared.HttpClient(true, url, string(out))
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}

	errChan <- &types.ImportError{Err: nil}
}

func createProducts(url string) error {

	var errChan = make(chan *types.ImportError)
	var wg sync.WaitGroup

	products, err := readProductsFile()
	if err != nil {
		shared.Error.Fatalln("Error reading file: ", err)
		return err
	}

	numProd := len(products)
	shared.Info.Printf("Found %d products in the file\n", numProd)
	shared.Info.Printf("Create products with %d connections\n", conn)

	if numProd < conn {
		wg.Add(numProd)
		for i := 0; i < numProd; i++ {
			go createAsyncProduct(url, products[i], &wg, errChan)
		}
		go func() {
			wg.Wait()
			close(errChan)
		}()
	} else {
		numOfLoops, remaining := numProd/conn, numProd%conn
		for i := 0; i < numOfLoops; i++ {
			shared.Info.Printf("Create %d batch of products\n", i)
			wg.Add(conn)
			for j := 0; j < conn; j++ {
				go createAsyncProduct(url, products[j+(i*conn)], &wg, errChan)
			}
			go func() {
				wg.Wait()
			}()
		}

		wg.Add(remaining)
		shared.Info.Printf("Create remaining %d products\n", remaining)
		for i := (numProd - remaining); i < numProd; i++ {
			go createAsyncProduct(url, products[i], &wg, errChan)
		}
		go func() {
			wg.Wait()
			close(errChan)
		}()
	}

	//print any errors and return an err
	var errs = false
	for errProd := range errChan {
		if errProd.Err != nil {
			shared.Error.Fatalln(errProd.Err)
			errs = true
		}
	}

	if errs {
		return fmt.Errorf("problem creating one of more products")
	}
	return nil
}

func readProductsFile() ([]Product, error) {

	products := []Product{}

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

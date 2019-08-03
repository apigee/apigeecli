package expprod

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	types "github.com/srinandan/apigeecli/cmd/types"
)

var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export API products to a file",
	Long:  "Export API products to a file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return exportProducts()
	},
}

type apiProducts struct {
	APIProducts []apiProduct `json:"apiProduct,omitempty"`
}

type apiProduct struct {
	Name string `json:"name,omitempty"`
}

var conn int

const file = "products.json"

func init() {

	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

}

func writeAsyncProduct(product string, wg *sync.WaitGroup, errChan chan<- *types.ImportError) {
	//this is a two step process - 1) get product details 2) write details to file
	defer wg.Done()

	u, _ := url.Parse(shared.BaseURL)
	u.Path = path.Join(u.Path, shared.RootArgs.Org, "apiproducts", product)
	//don't print to sysout
	respBody, err := shared.HttpClient(false, u.String())
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}

	defer f.Close()
	_, err = f.Write(respBody)
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}

	errChan <- &types.ImportError{Err: nil}
}

func exportProducts() error {

	var errChan = make(chan *types.ImportError)
	var wg sync.WaitGroup

	u, _ := url.Parse(shared.BaseURL)
	u.Path = path.Join(u.Path, shared.RootArgs.Org, "apiproducts")
	//don't print to sysout
	respBody, err := shared.HttpClient(false, u.String())
	if err != nil {
		return err
	}

	var products = apiProducts{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		return err
	}

	numProd := len(products.APIProducts)
	shared.Info.Printf("Found %d products in the org\n", numProd)
	shared.Info.Printf("Importing products with %d connections\n", conn)

	if numProd < conn {
		wg.Add(numProd)
		for i := 0; i < numProd; i++ {
			go writeAsyncProduct(url.PathEscape(products.APIProducts[i].Name), &wg, errChan)
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
		return fmt.Errorf("problem exporting one of more products")
	}
	return nil
}

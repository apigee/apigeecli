package expprod

import (
	"encoding/json"
	"fmt"
	"net/url"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		const exportFileName = "products.json"
		err = exportProducts()
		if err != nil {
			return
		}
		return shared.WriteJSONArrayToFile(exportFileName)
	},
}

type apiProducts struct {
	APIProducts []apiProduct `json:"apiProduct,omitempty"`
}

type apiProduct struct {
	Name string `json:"name,omitempty"`
}

var conn int

func init() {

	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

}

func exportProducts() error {

	var errChan = make(chan *types.ImportError)
	//parent workgroup
	var wg sync.WaitGroup
	var mu sync.Mutex
	const entityType = "apiproducts"

	u, _ := url.Parse(shared.BaseURL)
	u.Path = path.Join(u.Path, shared.RootArgs.Org, entityType)
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
	shared.Info.Printf("Exporting products with %d connections\n", conn)

	if numProd < conn {
		wg.Add(numProd)
		for i := 0; i < numProd; i++ {
			go shared.GetAsyncEntity(url.PathEscape(products.APIProducts[i].Name), entityType, &wg, &mu, errChan)
		}
		go func() {
			wg.Wait()
			close(errChan)
		}()
	} else {
		numOfLoops, remaining := numProd/conn, numProd%conn
		for i := 0; i < numOfLoops; i++ {
			shared.Info.Printf("Query %d batch of %d products\n", i+1, numProd)
			wg.Add(conn)
			for j := 0; j < conn; j++ {
				go shared.GetAsyncEntity(url.PathEscape(products.APIProducts[j+(i*conn)].Name), entityType, &wg, &mu, errChan)
			}
			go func() {
				wg.Wait()
			}()
		}

		wg.Add(remaining)
		shared.Info.Printf("Query remaining %d products\n", remaining)
		for i := (numProd - remaining); i < numProd; i++ {
			go shared.GetAsyncEntity(url.PathEscape(products.APIProducts[i].Name), entityType, &wg, &mu, errChan)
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

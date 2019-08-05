package expprod

import (
	"encoding/json"
	"net/url"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
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

//batch created a batch of products to query
func batch(apiProductNames []apiProduct, entityType string, pwg *sync.WaitGroup, mu *sync.Mutex) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(apiProductNames))

	for _, apiProductName := range apiProductNames {
		go shared.GetAsyncEntity(url.PathEscape(apiProductName.Name), entityType, &bwg, mu)
	}
	bwg.Wait()
}

func exportProducts() error {

	//parent workgroup
	var pwg sync.WaitGroup
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

	numOfLoops, remaining := numProd/conn, numProd%conn

	//ensure connections aren't greater then products
	if conn > numProd {
		conn = numProd
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		shared.Info.Printf("Exporting batch %d of products\n", (i + 1))
		go batch(products.APIProducts[start:end], entityType, &pwg, &mu)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		shared.Info.Printf("Exporting remaining %d products\n", remaining)
		go batch(products.APIProducts[start:numProd], entityType, &pwg, &mu)
		pwg.Wait()
	}

	return nil
}

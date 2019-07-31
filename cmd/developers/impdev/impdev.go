package impdev

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

type Developer struct {
	EMail      string            `json:"email,omitempty"`
	FirstName  string            `json:"firstName,omitempty"`
	LastName   string            `json:"lastName,omitempty"`
	Attributes []types.Attribute `json:"attributes,omitempty"`
	Username   string            `json:"userName,omitempty"`
}

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

func createAsyncDeveloper(url string, developer Developer, wg *sync.WaitGroup, errChan chan<- *types.ImportError) {
	defer wg.Done()
	out, err := json.Marshal(developer)
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}
	err = shared.HttpClient(url, string(out))
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}

	errChan <- &types.ImportError{Err: nil}
}

func createDevelopers(url string) error {

	var errChan = make(chan *types.ImportError)
	var wg sync.WaitGroup

	developers, err := readDevelopersFile()
	if err != nil {
		shared.Error.Fatalln("Error reading file: ", err)
		return err
	}

	numDev := len(developers)
	shared.Info.Printf("Found %d products in the file\n", numDev)
	shared.Info.Printf("Create products with %d connections\n", conn)

	if numDev < conn {
		wg.Add(numDev)
		for i := 0; i < numDev; i++ {
			go createAsyncDeveloper(url, developers[i], &wg, errChan)
		}
		go func() {
			wg.Wait()
			close(errChan)
		}()
	} else {
		numOfLoops, remaining := numDev/conn, numDev%conn
		for i := 0; i < numOfLoops; i++ {
			shared.Info.Printf("Create %d batch of products\n", i)
			wg.Add(conn)
			for j := 0; j < conn; j++ {
				go createAsyncDeveloper(url, developers[j+(i*conn)], &wg, errChan)
			}
			go func() {
				wg.Wait()
				close(errChan)
			}()
		}

		wg.Add(remaining)
		shared.Info.Printf("Create remaining %d products\n", remaining)
		for i := (numDev - remaining); i < numDev; i++ {
			go createAsyncDeveloper(url, developers[i], &wg, errChan)
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

func readDevelopersFile() ([]Developer, error) {

	developers := []Developer{}

	jsonFile, err := os.Open(file)

	if err != nil {
		return developers, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return developers, err
	}

	err = json.Unmarshal(byteValue, &developers)

	if err != nil {
		return developers, err
	}

	return developers, nil

}

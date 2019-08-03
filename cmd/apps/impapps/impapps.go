package impapps

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

type App struct {
	Name        string            `json:"name,omitempty"`
	Status      string            `json:"status,omitempty"`
	Credentials *[]Credential     `json:"credentials,omitempty"`
	DeveloperID *string           `json:"developerId,omitempty"`
	DisplayName string            `json:"displayName,omitempty"`
	Attributes  []types.Attribute `json:"attributes,omitempty"`
	CallbackURL string            `json:"callbackUrl,omitempty"`
	Scopes      []string          `json:"scopes,omitempty"`
}

type Credential struct {
	APIProducts    []APIProduct `json:"apiProducts,omitempty"`
	ConsumerKey    string       `json:"consumerKey,omitempty"`
	ConsumerSecret string       `json:"consumerSecret,omitempty"`
	ExpiresAt      int          `json:"expiresAt,omitempty"`
	Status         string       `json:"status,omitempty"`
	Scopes         []string     `json:"scopes,omitempty"`
}

type APIProduct struct {
	Name string `json:"apiproduct,omitempty"`
}

type ImportCredential struct {
	APIProducts    []string `json:"apiProducts,omitempty"`
	ConsumerKey    string   `json:"consumerKey,omitempty"`
	ConsumerSecret string   `json:"consumerSecret,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
}

var Cmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing Developer Apps",
	Long:  "Import a file containing Developer Apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		return createApps(shared.BaseURL)
	},
}

var conn int
var file string

func init() {

	Cmd.Flags().StringVarP(&file, "file", "f",
		"", "File containing Developer Apps")
	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = Cmd.MarkFlagRequired("file")
}

func createAsyncApp(app App, wg *sync.WaitGroup, errChan chan<- *types.ImportError) {
	defer wg.Done()

	//importing an app will be a two step process.
	//1. create the app without the credential
	//2. create/import the credential
	u, _ := url.Parse(shared.BaseURL)
	//store the developer and the credential
	developerID := *app.DeveloperID
	credentials := *app.Credentials

	//remove the developer id and credentials from the payload
	app.DeveloperID = nil
	app.Credentials = nil

	out, err := json.Marshal(app)
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}

	u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", developerID, "apps")
	_, err = shared.HttpClient(true, u.String(), string(out))
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}
	u, _ = url.Parse(shared.BaseURL)
	u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", developerID, "apps", app.Name, "keys", "create")
	for _, credential := range credentials {
		//construct a []string for products
		var products []string
		for _, apiProduct := range credential.APIProducts {
			products = append(products, apiProduct.Name)
		}
		//create a new credential
		importCredential := ImportCredential{}
		importCredential.APIProducts = products
		importCredential.ConsumerKey = credential.ConsumerKey
		importCredential.ConsumerSecret = credential.ConsumerSecret
		importCredential.Scopes = credential.Scopes

		impCred, err := json.Marshal(importCredential)
		if err != nil {
			errChan <- &types.ImportError{Err: err}
			return
		}
		_, err = shared.HttpClient(true, u.String(), string(impCred))
		if err != nil {
			errChan <- &types.ImportError{Err: err}
			return
		}
		shared.Warning.Println("NOTE: apiProducts are not associated with the app")
	}
	errChan <- &types.ImportError{Err: nil}
}

func createApps(url string) error {

	var errChan = make(chan *types.ImportError)
	var wg sync.WaitGroup

	apps, err := readAppsFile()
	if err != nil {
		shared.Error.Fatalln("Error reading file: ", err)
		return err
	}

	numApp := len(apps)
	shared.Info.Printf("Found %d apps in the file\n", numApp)
	shared.Info.Printf("Create apps with %d connections\n", conn)

	if numApp < conn {
		wg.Add(numApp)
		for i := 0; i < numApp; i++ {
			go createAsyncApp(apps[i], &wg, errChan)
		}
		go func() {
			wg.Wait()
			close(errChan)
		}()
	} else {
		numOfLoops, remaining := numApp/conn, numApp%conn
		for i := 0; i < numOfLoops; i++ {
			shared.Info.Printf("Create %d batch of apps\n", i)
			wg.Add(conn)
			for j := 0; j < conn; j++ {
				go createAsyncApp(apps[j+(i*conn)], &wg, errChan)
			}
			go func() {
				wg.Wait()
				close(errChan)
			}()
		}

		wg.Add(remaining)
		shared.Info.Printf("Create remaining %d apps\n", remaining)
		for i := (numApp - remaining); i < numApp; i++ {
			go createAsyncApp(apps[i], &wg, errChan)
		}
		go func() {
			wg.Wait()
			close(errChan)
		}()
	}

	//print any errors and return an err
	var errs = false
	for errApp := range errChan {
		if errApp.Err != nil {
			shared.Error.Fatalln(errApp.Err)
			errs = true
		}
	}

	if errs {
		return fmt.Errorf("problem creating one of more apps")
	}
	return nil
}

func readAppsFile() ([]App, error) {

	apps := []App{}

	jsonFile, err := os.Open(file)

	if err != nil {
		return apps, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return apps, err
	}

	err = json.Unmarshal(byteValue, &apps)

	if err != nil {
		return apps, err
	}

	return apps, nil

}

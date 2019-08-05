package impdev

import (
	"encoding/json"
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

func createAsyncDeveloper(url string, developer Developer, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := json.Marshal(developer)
	if err != nil {
		shared.Error.Fatalln(err)
		return
	}
	_, err = shared.HttpClient(true, url, string(out))
	if err != nil {
		shared.Error.Fatalln(err)
		return
	}

	shared.Info.Printf("Completed entity: %s", developer.EMail)
}

//batch creates a batch of developers to create
func batch(url string, entities []Developer, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		go createAsyncDeveloper(url, entity, &bwg)
	}
	bwg.Wait()
}

func createDevelopers(url string) error {

	var pwg sync.WaitGroup

	entities, err := readDevelopersFile()
	if err != nil {
		shared.Error.Fatalln("Error reading file: ", err)
		return err
	}

	numEntities := len(entities)
	shared.Info.Printf("Found %d developers in the file\n", numEntities)
	shared.Info.Printf("Create developers with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than products
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		shared.Info.Printf("Creating batch %d of developers\n", (i + 1))
		go batch(url, entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		shared.Info.Printf("Creating remaining %d developers\n", remaining)
		go batch(url, entities[start:numEntities], &pwg)
		pwg.Wait()
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

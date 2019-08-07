package impsfs

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "import",
	Short: "Import a folder containing sharedflow bundles",
	Long:  "Import a folder containing sharedflow bundles",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return createAPIs()
	},
}

var folder string
var conn int

func init() {

	Cmd.Flags().StringVarP(&folder, "folder", "f",
		"", "folder containing sharedflow bundles")
	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = Cmd.MarkFlagRequired("folder")
}

//batch creates a batch of sharedflow to import
func batch(entities []string, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		//sharedflow name is empty; same as filename
		go shared.ImportBundleAsync("sharedflows", "", entity, &bwg)
	}
	bwg.Wait()
}

func createAPIs() error {

	var pwg sync.WaitGroup
	var entities []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".zip" {
			return nil
		}
		entities = append(entities, path)
		return nil
	})

	if err != nil {
		return err
	}

	numEntities := len(entities)
	shared.Info.Printf("Found %d proxy bundles in the folder\n", numEntities)
	shared.Info.Printf("Create proxies with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		shared.Info.Printf("Creating batch %d of bundles\n", (i + 1))
		go batch(entities[start:end], &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		shared.Info.Printf("Creating remaining %d bundles\n", remaining)
		go batch(entities[start:numEntities], &pwg)
		pwg.Wait()
	}

	return nil
}

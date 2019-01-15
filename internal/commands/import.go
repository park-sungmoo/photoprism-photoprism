package commands

import (
	"fmt"

	"github.com/photoprism/photoprism/internal/context"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/urfave/cli"
)

// Imports photos from path defined in command-line args
var ImportCommand = cli.Command{
	Name:   "import",
	Usage:  "Imports photos",
	Action: importAction,
}

func importAction(ctx *cli.Context) error {
	conf := context.NewConfig(ctx)

	if err := conf.CreateDirectories(); err != nil {
		return err
	}

	conf.MigrateDb()

	fmt.Printf("Importing photos from %s...\n", conf.ImportPath())

	tensorFlow := photoprism.NewTensorFlow(conf.TensorFlowModelPath())

	indexer := photoprism.NewIndexer(conf.OriginalsPath(), tensorFlow, conf.Db())

	converter := photoprism.NewConverter(conf.DarktableCli())

	importer := photoprism.NewImporter(conf.OriginalsPath(), indexer, converter)

	importer.ImportPhotosFromDirectory(conf.ImportPath())

	fmt.Println("Done.")

	return nil
}

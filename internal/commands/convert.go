package commands

import (
	"fmt"
	"log"

	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/urfave/cli"
)

var ConvertCommand = cli.Command{
	Name:   "convert",
	Usage:  "Converts RAW originals to JPEG",
	Action: convertAction,
}

// Converts images to JPEG; called by ConvertCommand
func convertAction(context *cli.Context) error {
	conf := photoprism.NewConfig(context)

	if err := conf.CreateDirectories(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Converting RAW images in %s to JPEG...\n", conf.OriginalsPath)

	converter := photoprism.NewConverter(conf.DarktableCli)

	converter.ConvertAll(conf.OriginalsPath)

	fmt.Println("Done.")

	return nil
}

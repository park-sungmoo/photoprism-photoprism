package commands

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/photoprism/photoprism/pkg/fs"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/photoprism/photoprism/pkg/txt"
	"github.com/urfave/cli"
)

// IndexCommand registers the index cli command.
var IndexCommand = cli.Command{
	Name:      "index",
	Usage:     "Indexes media files in originals folder",
	UsageText: `To limit scope, a sub folder may be passed as first argument.`,
	Flags:     indexFlags,
	Action:    indexAction,
}

var indexFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "all, a",
		Usage: "re-index all originals, including unchanged files",
	},
	cli.BoolFlag{
		Name:  "cleanup",
		Usage: "removes orphan index entries and thumbnails",
	},
}

// indexAction indexes all photos in originals directory (photo library)
func indexAction(ctx *cli.Context) error {
	start := time.Now()

	conf := config.NewConfig(ctx)
	service.SetConfig(conf)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := conf.Init(); err != nil {
		return err
	}

	conf.InitDb()

	// Use first argument to limit scope if set.
	subPath := strings.TrimSpace(ctx.Args().First())

	if subPath == "" {
		log.Infof("indexing originals in %s", txt.Quote(conf.OriginalsPath()))
	} else {
		log.Infof("indexing originals in %s", txt.Quote(filepath.Join(conf.OriginalsPath(), subPath)))
	}

	if conf.ReadOnly() {
		log.Infof("index: read-only mode enabled")
	}

	var indexed fs.Done

	if w := service.Index(); w != nil {
		opt := photoprism.IndexOptions{
			Path:    subPath,
			Rescan:  ctx.Bool("all"),
			Convert: conf.Settings().Index.Convert && conf.SidecarWritable(),
			Stack:   true,
		}

		indexed = w.Start(opt)
	}

	if w := service.Purge(); w != nil {
		opt := photoprism.PurgeOptions{
			Path:   subPath,
			Ignore: indexed,
		}

		if files, photos, err := w.Start(opt); err != nil {
			log.Error(err)
		} else if len(files) > 0 || len(photos) > 0 {
			log.Infof("purge: removed %d files and %d photos", len(files), len(photos))
		}
	}

	if ctx.Bool("cleanup") {
		w := service.CleanUp()

		opt := photoprism.CleanUpOptions{
			Dry: false,
		}

		if thumbs, orphans, err := w.Start(opt); err != nil {
			return err
		} else {
			log.Infof("cleanup: removed %d index entries and %d orphan thumbnails", orphans, thumbs)
		}
	}

	elapsed := time.Since(start)

	log.Infof("indexed %d files in %s", len(indexed), elapsed)

	conf.Shutdown()

	return nil
}

package commands

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize/english"

	"github.com/urfave/cli"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/txt"
)

// IndexCommand registers the index cli command.
var IndexCommand = cli.Command{
	Name:      "index",
	Usage:     "Indexes original media files",
	ArgsUsage: "[path]",
	Flags:     indexFlags,
	Action:    indexAction,
}

var indexFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "force, f",
		Usage: "re-index all originals, including unchanged files",
	},
	cli.BoolFlag{
		Name:  "cleanup, c",
		Usage: "remove orphan index entries and thumbnails",
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
			Rescan:  ctx.Bool("force"),
			Convert: conf.Settings().Index.Convert && conf.SidecarWritable(),
			Stack:   true,
		}

		indexed = w.Start(opt)
	}

	if w := service.Purge(); w != nil {
		purgeStart := time.Now()
		opt := photoprism.PurgeOptions{
			Path:   subPath,
			Ignore: indexed,
		}

		if files, photos, err := w.Start(opt); err != nil {
			log.Error(err)
		} else if len(files) > 0 || len(photos) > 0 {
			log.Infof("purge: removed %s and %s [%s]", english.Plural(len(files), "file", "files"), english.Plural(len(photos), "photo", "photos"), time.Since(purgeStart))
		}
	}

	if ctx.Bool("cleanup") {
		cleanupStart := time.Now()
		w := service.CleanUp()

		opt := photoprism.CleanUpOptions{
			Dry: false,
		}

		if thumbs, orphans, err := w.Start(opt); err != nil {
			return err
		} else {
			log.Infof("cleanup: removed %s and %s [%s]", english.Plural(orphans, "index entry", "index entries"), english.Plural(thumbs, "thumbnail", "thumbnails"), time.Since(cleanupStart))
		}
	}

	elapsed := time.Since(start)

	log.Infof("indexed %s in %s", english.Plural(len(indexed), "file", "files"), elapsed)

	conf.Shutdown()

	return nil
}

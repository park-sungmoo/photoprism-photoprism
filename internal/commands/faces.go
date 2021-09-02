package commands

import (
	"context"
	"time"

	"github.com/manifoldco/promptui"

	"github.com/photoprism/photoprism/internal/photoprism"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/urfave/cli"
)

// FacesCommand registers the faces cli command.
var FacesCommand = cli.Command{
	Name:  "faces",
	Usage: "Facial recognition sub-commands",
	Subcommands: []cli.Command{
		{
			Name:   "stats",
			Usage:  "Shows stats on face samples",
			Action: facesStatsAction,
		},
		{
			Name:  "audit",
			Usage: "Conducts a data integrity audit",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "fix, f",
					Usage: "issues will be fixed automatically",
				},
			},
			Action: facesAuditAction,
		},
		{
			Name:   "reset",
			Usage:  "Resets recognized faces",
			Action: facesResetAction,
		},
		{
			Name:   "optimize",
			Usage:  "Optimizes face clusters",
			Action: facesOptimizeAction,
		},
		{
			Name:  "update",
			Usage: "Performs facial recognition",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "update existing faces",
				},
			},
			Action: facesUpdateAction,
		},
	},
}

// facesStatsAction shows stats on face embeddings.
func facesStatsAction(ctx *cli.Context) error {
	start := time.Now()

	conf := config.NewConfig(ctx)
	service.SetConfig(conf)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := conf.Init(); err != nil {
		return err
	}

	conf.InitDb()

	w := service.Faces()

	if err := w.Stats(); err != nil {
		return err
	} else {
		elapsed := time.Since(start)

		log.Infof("completed in %s", elapsed)
	}

	conf.Shutdown()

	return nil
}

// facesAuditAction shows stats on face embeddings.
func facesAuditAction(ctx *cli.Context) error {
	start := time.Now()

	conf := config.NewConfig(ctx)
	service.SetConfig(conf)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := conf.Init(); err != nil {
		return err
	}

	conf.InitDb()

	w := service.Faces()

	if err := w.Audit(ctx.Bool("fix")); err != nil {
		return err
	} else {
		elapsed := time.Since(start)

		log.Infof("completed in %s", elapsed)
	}

	conf.Shutdown()

	return nil
}

// facesResetAction resets face clusters and matches.
func facesResetAction(ctx *cli.Context) error {
	actionPrompt := promptui.Prompt{
		Label:     "Remove automatically recognized faces, matches, and dangling subjects?",
		IsConfirm: true,
	}

	if _, err := actionPrompt.Run(); err != nil {
		return nil
	}

	start := time.Now()

	conf := config.NewConfig(ctx)
	service.SetConfig(conf)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := conf.Init(); err != nil {
		return err
	}

	conf.InitDb()

	w := service.Faces()

	if err := w.Reset(); err != nil {
		return err
	} else {
		elapsed := time.Since(start)

		log.Infof("completed in %s", elapsed)
	}

	conf.Shutdown()

	return nil
}

// facesOptimizeAction optimizes existing face clusters.
func facesOptimizeAction(ctx *cli.Context) error {
	start := time.Now()

	conf := config.NewConfig(ctx)
	service.SetConfig(conf)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := conf.Init(); err != nil {
		return err
	}

	conf.InitDb()

	w := service.Faces()

	if res, err := w.Optimize(); err != nil {
		return err
	} else {
		elapsed := time.Since(start)

		log.Infof("%d face clusters merged in %s", res.Merged, elapsed)
	}

	conf.Shutdown()

	return nil
}

// facesUpdateAction performs face clustering and matching.
func facesUpdateAction(ctx *cli.Context) error {
	start := time.Now()

	conf := config.NewConfig(ctx)
	service.SetConfig(conf)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := conf.Init(); err != nil {
		return err
	}

	conf.InitDb()

	opt := photoprism.FacesOptions{
		Force: ctx.Bool("force"),
	}

	w := service.Faces()

	if err := w.Start(opt); err != nil {
		return err
	} else {
		elapsed := time.Since(start)

		log.Infof("completed in %s", elapsed)
	}

	conf.Shutdown()

	return nil
}

package commands

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/urfave/cli"
)

// ConfigCommand is used to register the display config cli command
var ConfigCommand = cli.Command{
	Name:   "config",
	Usage:  "Displays global configuration values",
	Action: configAction,
}

// configAction prints current configuration
func configAction(ctx *cli.Context) error {
	conf := config.NewConfig(ctx)

	dbDriver := conf.DatabaseDriver()
	dbDsn := conf.DatabaseDsn()

	fmt.Printf("%-25s VALUE\n", "NAME")

	// Feature flags.
	fmt.Printf("%-25s %t\n", "debug", conf.Debug())
	fmt.Printf("%-25s %t\n", "public", conf.Public())
	fmt.Printf("%-25s %t\n", "read-only", conf.ReadOnly())
	fmt.Printf("%-25s %t\n", "experimental", conf.Experimental())

	// Config path and main file.
	fmt.Printf("%-25s %s\n", "config-path", conf.ConfigPath())
	fmt.Printf("%-25s %s\n", "config-file", conf.ConfigFile())
	fmt.Printf("%-25s %s\n", "settings-file", conf.SettingsFile())

	// Main directories.
	fmt.Printf("%-25s %s\n", "originals-path", conf.OriginalsPath())
	fmt.Printf("%-25s %d\n", "originals-limit", conf.OriginalsLimit())
	fmt.Printf("%-25s %s\n", "import-path", conf.ImportPath())
	fmt.Printf("%-25s %s\n", "storage-path", conf.StoragePath())
	fmt.Printf("%-25s %s\n", "sidecar-path", conf.SidecarPath())
	fmt.Printf("%-25s %s\n", "albums-path", conf.AlbumsPath())
	fmt.Printf("%-25s %s\n", "cache-path", conf.CachePath())
	fmt.Printf("%-25s %s\n", "temp-path", conf.TempPath())
	fmt.Printf("%-25s %s\n", "backup-path", conf.BackupPath())
	fmt.Printf("%-25s %s\n", "assets-path", conf.AssetsPath())

	// Asset path and file names.
	fmt.Printf("%-25s %s\n", "static-path", conf.StaticPath())
	fmt.Printf("%-25s %s\n", "build-path", conf.BuildPath())
	fmt.Printf("%-25s %s\n", "img-path", conf.ImgPath())
	fmt.Printf("%-25s %s\n", "templates-path", conf.TemplatesPath())

	// Workers.
	fmt.Printf("%-25s %d\n", "workers", conf.Workers())
	fmt.Printf("%-25s %d\n", "wakeup-interval", conf.WakeupInterval()/time.Second)

	// Disable features.
	fmt.Printf("%-25s %t\n", "disable-backups", conf.DisableBackups())
	fmt.Printf("%-25s %t\n", "disable-settings", conf.DisableSettings())
	fmt.Printf("%-25s %t\n", "disable-places", conf.DisablePlaces())
	fmt.Printf("%-25s %t\n", "disable-exiftool", conf.DisableExifTool())

	// Everything related to TensorFlow.
	fmt.Printf("%-25s %t\n", "disable-tensorflow", conf.DisableTensorFlow())
	fmt.Printf("%-25s %s\n", "tensorflow-version", conf.TensorFlowVersion())
	fmt.Printf("%-25s %s\n", "tensorflow-model-path", conf.TensorFlowModelPath())
	fmt.Printf("%-25s %t\n", "detect-nsfw", conf.DetectNSFW())
	fmt.Printf("%-25s %t\n", "upload-nsfw", conf.UploadNSFW())

	// Site information.
	fmt.Printf("%-25s %s\n", "site-url", conf.SiteUrl())
	fmt.Printf("%-25s %s\n", "site-preview", conf.SitePreview())
	fmt.Printf("%-25s %s\n", "site-title", conf.SiteTitle())
	fmt.Printf("%-25s %s\n", "site-caption", conf.SiteCaption())
	fmt.Printf("%-25s %s\n", "site-description", conf.SiteDescription())
	fmt.Printf("%-25s %s\n", "site-author", conf.SiteAuthor())

	// Logging.
	fmt.Printf("%-25s %s\n", "log-level", conf.LogLevel())
	fmt.Printf("%-25s %s\n", "log-filename", conf.LogFilename())
	fmt.Printf("%-25s %s\n", "pid-filename", conf.PIDFilename())

	// HTTP server configuration.
	fmt.Printf("%-25s %s\n", "http-host", conf.HttpHost())
	fmt.Printf("%-25s %d\n", "http-port", conf.HttpPort())
	fmt.Printf("%-25s %s\n", "http-mode", conf.HttpMode())

	// Passwords.
	fmt.Printf("%-25s %s\n", "admin-password", strings.Repeat("*", utf8.RuneCountInString(conf.AdminPassword())))

	// Database configuration.
	fmt.Printf("%-25s %s\n", "database-driver", dbDriver)
	fmt.Printf("%-25s %s\n", "database-dsn", dbDsn)
	fmt.Printf("%-25s %s\n", "database-server", conf.DatabaseServer())
	fmt.Printf("%-25s %s\n", "database-host", conf.DatabaseHost())
	fmt.Printf("%-25s %s\n", "database-port", conf.DatabasePortString())
	fmt.Printf("%-25s %s\n", "database-name", conf.DatabaseName())
	fmt.Printf("%-25s %s\n", "database-user", conf.DatabaseUser())
	fmt.Printf("%-25s %s\n", "database-password", strings.Repeat("*", utf8.RuneCountInString(conf.DatabasePassword())))
	fmt.Printf("%-25s %d\n", "database-conns", conf.DatabaseConns())
	fmt.Printf("%-25s %d\n", "database-conns-idle", conf.DatabaseConnsIdle())

	// External binaries and sidecar configuration.
	fmt.Printf("%-25s %s\n", "rawtherapee-bin", conf.RawtherapeeBin())
	fmt.Printf("%-25s %s\n", "darktable-bin", conf.DarktableBin())
	fmt.Printf("%-25s %t\n", "darktable-presets", conf.DarktablePresets())
	fmt.Printf("%-25s %s\n", "sips-bin", conf.SipsBin())
	fmt.Printf("%-25s %s\n", "heifconvert-bin", conf.HeifConvertBin())
	fmt.Printf("%-25s %s\n", "ffmpeg-bin", conf.FFmpegBin())
	fmt.Printf("%-25s %s\n", "exiftool-bin", conf.ExifToolBin())

	// Thumbs, resampling and download security token.
	fmt.Printf("%-25s %s\n", "download-token", conf.DownloadToken())
	fmt.Printf("%-25s %s\n", "preview-token", conf.PreviewToken())
	fmt.Printf("%-25s %s\n", "thumb-filter", conf.ThumbFilter())
	fmt.Printf("%-25s %t\n", "thumb-uncached", conf.ThumbUncached())
	fmt.Printf("%-25s %d\n", "thumb-size", conf.ThumbSize())
	fmt.Printf("%-25s %d\n", "thumb-size-uncached", conf.ThumbSizeUncached())
	fmt.Printf("%-25s %s\n", "thumb-path", conf.ThumbPath())
	fmt.Printf("%-25s %d\n", "jpeg-size", conf.JpegSize())
	fmt.Printf("%-25s %d\n", "jpeg-quality", conf.JpegQuality())

	return nil
}

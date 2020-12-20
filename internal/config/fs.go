package config

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/txt"
)

func findExecutable(configBin, defaultBin string) (result string) {
	if configBin == "" {
		result = defaultBin
	} else {
		result = configBin
	}

	if path, err := exec.LookPath(result); err == nil {
		result = path
	}

	if !fs.FileExists(result) {
		result = ""
	}

	return result
}

// CreateDirectories creates directories for storing photos, metadata and cache files.
func (c *Config) CreateDirectories() error {
	createError := func(path string, err error) (result error) {
		if fs.FileExists(path) {
			result = fmt.Errorf("%s is a file, not a folder: please check your configuration", txt.Quote(path))
		} else {
			result = fmt.Errorf("can't create %s: please check configuration and permissions", txt.Quote(path))
		}

		log.Debug(err)

		return result
	}

	notFoundError := func(name string) error {
		return fmt.Errorf("%s path not found, run 'photoprism config' to check configuration options", name)
	}

	if c.AssetsPath() == "" {
		return notFoundError("assets")
	} else if err := os.MkdirAll(c.AssetsPath(), os.ModePerm); err != nil {
		return createError(c.AssetsPath(), err)
	}

	if c.StoragePath() == "" {
		return notFoundError("storage")
	} else if err := os.MkdirAll(c.StoragePath(), os.ModePerm); err != nil {
		return createError(c.StoragePath(), err)
	}

	if c.BackupPath() == "" {
		return notFoundError("backup")
	} else if err := os.MkdirAll(c.BackupPath(), os.ModePerm); err != nil {
		return createError(c.BackupPath(), err)
	}

	if c.OriginalsPath() == "" {
		return notFoundError("originals")
	} else if err := os.MkdirAll(c.OriginalsPath(), os.ModePerm); err != nil {
		return createError(c.OriginalsPath(), err)
	}

	if c.ImportPath() == "" {
		return notFoundError("import")
	} else if err := os.MkdirAll(c.ImportPath(), os.ModePerm); err != nil {
		return createError(c.ImportPath(), err)
	}

	if filepath.IsAbs(c.SidecarPath()) {
		if err := os.MkdirAll(c.SidecarPath(), os.ModePerm); err != nil {
			return createError(c.SidecarPath(), err)
		}
	}

	if c.CachePath() == "" {
		return notFoundError("cache")
	} else if err := os.MkdirAll(c.CachePath(), os.ModePerm); err != nil {
		return createError(c.CachePath(), err)
	}

	if c.ThumbPath() == "" {
		return notFoundError("thumbs")
	} else if err := os.MkdirAll(c.ThumbPath(), os.ModePerm); err != nil {
		return createError(c.ThumbPath(), err)
	}

	if c.ConfigPath() == "" {
		return notFoundError("config")
	} else if err := os.MkdirAll(c.ConfigPath(), os.ModePerm); err != nil {
		return createError(c.ConfigPath(), err)
	}

	if c.TempPath() == "" {
		return notFoundError("temp")
	} else if err := os.MkdirAll(c.TempPath(), os.ModePerm); err != nil {
		return createError(c.TempPath(), err)
	}

	if c.AlbumsPath() == "" {
		return notFoundError("albums")
	} else if err := os.MkdirAll(c.AlbumsPath(), os.ModePerm); err != nil {
		return createError(c.AlbumsPath(), err)
	}

	if c.TensorFlowModelPath() == "" {
		return notFoundError("tensorflow model")
	} else if err := os.MkdirAll(c.TensorFlowModelPath(), os.ModePerm); err != nil {
		return createError(c.TensorFlowModelPath(), err)
	}

	if c.BuildPath() == "" {
		return notFoundError("build")
	} else if err := os.MkdirAll(c.BuildPath(), os.ModePerm); err != nil {
		return createError(c.BuildPath(), err)
	}

	if filepath.Dir(c.PIDFilename()) == "" {
		return notFoundError("pid file")
	} else if err := os.MkdirAll(filepath.Dir(c.PIDFilename()), os.ModePerm); err != nil {
		return createError(filepath.Dir(c.PIDFilename()), err)
	}

	if filepath.Dir(c.LogFilename()) == "" {
		return notFoundError("log file")
	} else if err := os.MkdirAll(filepath.Dir(c.LogFilename()), os.ModePerm); err != nil {
		return createError(filepath.Dir(c.LogFilename()), err)
	}

	return nil
}

// ConfigFile returns the config file name.
func (c *Config) ConfigFile() string {
	if c.options.ConfigFile == "" || !fs.FileExists(c.options.ConfigFile) {
		return filepath.Join(c.ConfigPath(), "options.yml")
	}

	return c.options.ConfigFile
}

// ConfigPath returns the config path.
func (c *Config) ConfigPath() string {
	if c.options.ConfigPath == "" {
		if fs.PathExists(filepath.Join(c.StoragePath(), "settings")) {
			return filepath.Join(c.StoragePath(), "settings")
		}

		return filepath.Join(c.StoragePath(), "config")
	}

	return fs.Abs(c.options.ConfigPath)
}

// HubConfigFile returns the backend api config file name.
func (c *Config) HubConfigFile() string {
	return filepath.Join(c.ConfigPath(), "hub.yml")
}

// SettingsFile returns the user settings file name.
func (c *Config) SettingsFile() string {
	return filepath.Join(c.ConfigPath(), "settings.yml")
}

// PIDFilename returns the filename for storing the server process id (pid).
func (c *Config) PIDFilename() string {
	if c.options.PIDFilename == "" {
		return filepath.Join(c.StoragePath(), "photoprism.pid")
	}

	return fs.Abs(c.options.PIDFilename)
}

// LogFilename returns the filename for storing server logs.
func (c *Config) LogFilename() string {
	if c.options.LogFilename == "" {
		return filepath.Join(c.StoragePath(), "photoprism.log")
	}

	return fs.Abs(c.options.LogFilename)
}

// OriginalsPath returns the originals.
func (c *Config) OriginalsPath() string {
	if c.options.OriginalsPath == "" {
		// Try to find the right directory by iterating through a list.
		c.options.OriginalsPath = fs.FindDir(fs.OriginalPaths)
	}

	return fs.Abs(c.options.OriginalsPath)
}

// ImportPath returns the import directory.
func (c *Config) ImportPath() string {
	if c.options.ImportPath == "" {
		// Try to find the right directory by iterating through a list.
		c.options.ImportPath = fs.FindDir(fs.ImportPaths)
	}

	return fs.Abs(c.options.ImportPath)
}

// ExifToolBin returns the exiftool executable file name.
func (c *Config) ExifToolBin() string {
	return findExecutable(c.options.ExifToolBin, "exiftool")
}

// Automatically create JSON sidecar files using Exiftool.
func (c *Config) ExifToolJson() bool {
	return !c.DisableExifTool()
}

// Automatically backup metadata to YAML sidecar files.
func (c *Config) BackupYaml() bool {
	return !c.DisableBackups()
}

// SidecarPath returns the storage path for generated sidecar files (relative or absolute).
func (c *Config) SidecarPath() string {
	if c.options.SidecarPath == "" {
		c.options.SidecarPath = filepath.Join(c.StoragePath(), "sidecar")
	}

	return c.options.SidecarPath
}

// SidecarPathIsAbs tests if sidecar path is absolute.
func (c *Config) SidecarPathIsAbs() bool {
	return filepath.IsAbs(c.SidecarPath())
}

// SidecarWritable tests if sidecar files can be created.
func (c *Config) SidecarWritable() bool {
	return !c.ReadOnly() || c.SidecarPathIsAbs()
}

// FFmpegBin returns the ffmpeg executable file name.
func (c *Config) FFmpegBin() string {
	return findExecutable(c.options.FFmpegBin, "ffmpeg")
}

// TempPath returns a temporary directory name for uploads and downloads.
func (c *Config) TempPath() string {
	if c.options.TempPath == "" {
		return filepath.Join(os.TempDir(), "photoprism")
	}

	return fs.Abs(c.options.TempPath)
}

// CachePath returns the path for cache files.
func (c *Config) CachePath() string {
	if c.options.CachePath == "" {
		return filepath.Join(c.StoragePath(), "cache")
	}

	return fs.Abs(c.options.CachePath)
}

// StoragePath returns the path for generated files like cache and index.
func (c *Config) StoragePath() string {
	if c.options.StoragePath == "" {
		const dirName = "storage"

		// Default directories.
		originalsDir := fs.Abs(filepath.Join(c.OriginalsPath(), fs.HiddenPath, dirName))
		storageDir := fs.Abs(dirName)

		// Find existing directories.
		if fs.PathWritable(originalsDir) && !c.ReadOnly() {
			return originalsDir
		} else if fs.PathWritable(storageDir) && c.ReadOnly() {
			return storageDir
		}

		// Fallback to backup storage path.
		if fs.PathWritable(c.options.BackupPath) {
			return fs.Abs(filepath.Join(c.options.BackupPath, dirName))
		}

		// Use .photoprism in home directory?
		if usr, _ := user.Current(); usr.HomeDir != "" {
			p := fs.Abs(filepath.Join(usr.HomeDir, fs.HiddenPath, dirName))

			if fs.PathWritable(p) || c.ReadOnly() {
				return p
			}
		}

		// Fallback directory in case nothing else works.
		if c.ReadOnly() {
			return fs.Abs(filepath.Join(fs.HiddenPath, dirName))
		}

		// Store cache and index in "originals/.photoprism/storage".
		return originalsDir
	}

	return fs.Abs(c.options.StoragePath)
}

// BackupPath returns the backup storage path.
func (c *Config) BackupPath() string {
	if fs.PathWritable(c.options.BackupPath) {
		return fs.Abs(c.options.BackupPath)
	}

	return filepath.Join(c.StoragePath(), "backup")
}

// AssetsPath returns the path to static assets for models and templates.
func (c *Config) AssetsPath() string {
	if c.options.AssetsPath == "" {
		// Try to find the right directory by iterating through a list.
		c.options.AssetsPath = fs.FindDir(fs.AssetPaths)
	}

	return fs.Abs(c.options.AssetsPath)
}

// LocalesPath returns the translation locales path.
func (c *Config) LocalesPath() string {
	return filepath.Join(c.AssetsPath(), "locales")
}

// ExamplesPath returns the example files path.
func (c *Config) ExamplesPath() string {
	return filepath.Join(c.AssetsPath(), "examples")
}

// TestdataPath returns the test files path.
func (c *Config) TestdataPath() string {
	return filepath.Join(c.StoragePath(), "testdata")
}

// MysqlBin returns the mysql executable file name.
func (c *Config) MysqlBin() string {
	return findExecutable("", "mysql")
}

// MysqldumpBin returns the mysqldump executable file name.
func (c *Config) MysqldumpBin() string {
	return findExecutable("", "mysqldump")
}

// SqliteBin returns the sqlite executable file name.
func (c *Config) SqliteBin() string {
	return findExecutable("", "sqlite3")
}

// AlbumsPath returns the storage path for album YAML files.
func (c *Config) AlbumsPath() string {
	return filepath.Join(c.StoragePath(), "albums")
}

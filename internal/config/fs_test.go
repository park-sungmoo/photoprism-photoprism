package config

import (
	"testing"

	"github.com/photoprism/photoprism/pkg/rnd"
	"github.com/stretchr/testify/assert"
)

func TestConfig_FindExecutable(t *testing.T) {
	assert.Equal(t, "", findExecutable("yyy", "xxx"))
}

func TestConfig_SidecarJson(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, true, c.ExifToolJson())
	assert.Equal(t, c.DisableExifTool(), !c.ExifToolJson())

	c.options.DisableExifTool = true

	assert.Equal(t, false, c.ExifToolJson())
	assert.Equal(t, c.DisableExifTool(), !c.ExifToolJson())
}

func TestConfig_SidecarYaml(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, true, c.BackupYaml())
	assert.Equal(t, c.DisableBackups(), !c.BackupYaml())

	c.options.DisableBackups = true

	assert.Equal(t, false, c.BackupYaml())
	assert.Equal(t, c.DisableBackups(), !c.BackupYaml())
}

func TestConfig_SidecarPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Contains(t, c.SidecarPath(), "testdata/sidecar")
	c.options.SidecarPath = ".photoprism"
	assert.Equal(t, ".photoprism", c.SidecarPath())
}

func TestConfig_SidecarPathIsAbs(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, true, c.SidecarPathIsAbs())
	c.options.SidecarPath = ".photoprism"
	assert.Equal(t, false, c.SidecarPathIsAbs())
}

func TestConfig_SidecarWritable(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, true, c.SidecarWritable())
}

func TestConfig_FFmpegBin(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "/usr/bin/ffmpeg", c.FFmpegBin())
}

func TestConfig_TempPath(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/temp", c.TempPath())
	c.options.TempPath = ""
	assert.Equal(t, "/tmp/photoprism", c.TempPath())
}

func TestConfig_CachePath2(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/cache", c.CachePath())
	c.options.CachePath = ""
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/cache", c.CachePath())
}

func TestConfig_StoragePath(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata", c.StoragePath())
	c.options.StoragePath = ""
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/originals/.photoprism/storage", c.StoragePath())
}

func TestConfig_TestdataPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/testdata", c.TestdataPath())
}

func TestConfig_CreateDirectories(t *testing.T) {
	testConfigMutex.Lock()
	defer testConfigMutex.Unlock()

	c := &Config{
		options: NewTestOptions(),
		token:   rnd.Token(8),
	}

	if err := c.CreateDirectories(); err != nil {
		t.Fatal(err)
	}
}

func TestConfig_ConfigFile2(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/config/options.yml", c.ConfigFile())
	c.options.ConfigFile = "/go/src/github.com/photoprism/photoprism/internal/config/testdata/config.yml"
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/internal/config/testdata/config.yml", c.ConfigFile())
}

func TestConfig_PIDFilename2(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/photoprism.pid", c.PIDFilename())
	c.options.PIDFilename = "/go/src/github.com/photoprism/photoprism/internal/config/testdata/test.pid"
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/internal/config/testdata/test.pid", c.PIDFilename())
}

func TestConfig_LogFilename2(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/photoprism.log", c.LogFilename())
	c.options.LogFilename = "/go/src/github.com/photoprism/photoprism/internal/config/testdata/test.log"
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/internal/config/testdata/test.log", c.LogFilename())
}

func TestConfig_OriginalsPath2(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/originals", c.OriginalsPath())
	c.options.OriginalsPath = ""
	assert.Equal(t, "", c.OriginalsPath())
}

func TestConfig_ImportPath2(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/storage/testdata/import", c.ImportPath())
	c.options.ImportPath = ""
	assert.Equal(t, "", c.ImportPath())
}

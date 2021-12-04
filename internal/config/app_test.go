package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_AppIcon(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "logo", c.AppIcon())
	c.options.AppIcon = "foo"
	assert.Equal(t, "logo", c.AppIcon())
	c.options.AppIcon = "app"
	assert.Equal(t, "app", c.AppIcon())
	c.options.AppIcon = "crisp"
	assert.Equal(t, "crisp", c.AppIcon())
	c.options.AppIcon = "mint"
	assert.Equal(t, "mint", c.AppIcon())
	c.options.AppIcon = "bold"
	assert.Equal(t, "bold", c.AppIcon())
	c.options.AppIcon = "logo"
	assert.Equal(t, "logo", c.AppIcon())
}

func TestConfig_AppIconsPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	if p := c.AppIconsPath(); !strings.HasSuffix(p, "photoprism/assets/static/icons") {
		t.Fatal("path .../photoprism/assets/static/icons expected")
	}

	if p := c.AppIconsPath("app"); !strings.HasSuffix(p, "photoprism/assets/static/icons/app") {
		t.Fatal("path .../pphotoprism/assets/static/icons/app expected")
	}

	if p := c.AppIconsPath("app", "512.png"); !strings.HasSuffix(p, "photoprism/assets/static/icons/app/512.png") {
		t.Fatal("path .../photoprism/assets/static/icons/app/512.png expected")
	}
}

func TestConfig_AppName(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "config.test", c.AppName())
}

func TestConfig_AppMode(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "standalone", c.AppMode())
}

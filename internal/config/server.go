package config

import (
	"path/filepath"

	"github.com/photoprism/photoprism/pkg/fs"
)

// DetachServer tests if server should detach from console (daemon mode).
func (c *Config) DetachServer() bool {
	return c.options.DetachServer
}

// HttpHost returns the built-in HTTP server host name or IP address (empty for all interfaces).
func (c *Config) HttpHost() string {
	if c.options.HttpHost == "" {
		return "0.0.0.0"
	}

	return c.options.HttpHost
}

// HttpPort returns the built-in HTTP server port.
func (c *Config) HttpPort() int {
	if c.options.HttpPort == 0 {
		return 2342
	}

	return c.options.HttpPort
}

// HttpMode returns the server mode.
func (c *Config) HttpMode() string {
	if c.options.HttpMode == "" {
		if c.Debug() {
			return "debug"
		}

		return "release"
	}

	return c.options.HttpMode
}

// TemplatesPath returns the server templates path.
func (c *Config) TemplatesPath() string {
	return filepath.Join(c.AssetsPath(), "templates")
}

// TemplateExists tests if a template with the given name exists (e.g. index.tmpl).
func (c *Config) TemplateExists(name string) bool {
	return fs.FileExists(filepath.Join(c.TemplatesPath(), name))
}

// TemplateName returns the name of the default template (e.g. index.tmpl).
func (c *Config) TemplateName() string {
	if s := c.Settings(); s != nil {
		if c.TemplateExists(s.Templates.Default) {
			return s.Templates.Default
		}
	}

	return "index.tmpl"
}

// StaticPath returns the static assets path.
func (c *Config) StaticPath() string {
	return filepath.Join(c.AssetsPath(), "static")
}

// BuildPath returns the static build path.
func (c *Config) BuildPath() string {
	return filepath.Join(c.StaticPath(), "build")
}

// ImgPath returns the static image path.
func (c *Config) ImgPath() string {
	return filepath.Join(c.StaticPath(), "img")
}

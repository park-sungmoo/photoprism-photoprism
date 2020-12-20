package config

import (
	"regexp"

	"github.com/photoprism/photoprism/pkg/rnd"
	"golang.org/x/crypto/bcrypt"
)

func isBcrypt(s string) bool {
	b, err := regexp.MatchString(`^\$2[ayb]\$.{56}$`, s)
	if err != nil {
		return false
	}
	return b
}

// CheckPassword compares given password p with the admin password
func (c *Config) CheckPassword(p string) bool {
	ap := c.AdminPassword()

	if isBcrypt(ap) {
		err := bcrypt.CompareHashAndPassword([]byte(ap), []byte(p))
		return err == nil
	}

	return ap == p
}

// InvalidDownloadToken tests if the token is invalid.
func (c *Config) InvalidDownloadToken(t string) bool {
	return c.DownloadToken() != t
}

// DownloadToken returns the DOWNLOAD api token (you can optionally use a static value for permanent caching).
func (c *Config) DownloadToken() string {
	if c.options.DownloadToken == "" {
		c.options.DownloadToken = rnd.Token(8)
	}

	return c.options.DownloadToken
}

// InvalidPreviewToken tests if the preview token is invalid.
func (c *Config) InvalidPreviewToken(t string) bool {
	return c.PreviewToken() != t && c.DownloadToken() != t
}

// PreviewToken returns the THUMBNAILS api token (you can optionally use a static value for permanent caching).
func (c *Config) PreviewToken() string {
	if c.options.PreviewToken == "" {
		c.options.PreviewToken = rnd.Token(8)
	}

	return c.options.PreviewToken
}

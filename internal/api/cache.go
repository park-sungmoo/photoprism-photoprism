package api

import (
	"fmt"
	"strconv"

	"github.com/photoprism/photoprism/internal/query"

	"github.com/photoprism/photoprism/internal/service"
	"github.com/photoprism/photoprism/internal/thumb"
)

// MaxAge represents a cache TTL in seconds.
type MaxAge int

// String returns the cache TTL in seconds as string.
func (a MaxAge) String() string {
	return strconv.Itoa(int(a))
}

// Default cache TTL times in seconds.
var (
	CoverCacheTTL MaxAge = 3600           // 1 hour
	ThumbCacheTTL MaxAge = 3600 * 24 * 90 // ~ 3 months
)

type ThumbCache struct {
	FileName  string
	ShareName string
}

type ByteCache struct {
	Data []byte
}

// CacheKey returns a cache key string based on namespace, uid and name.
func CacheKey(ns, uid, name string) string {
	return fmt.Sprintf("%s:%s:%s", ns, uid, name)
}

// RemoveFromFolderCache removes an item from the folder cache e.g. after indexing.
func RemoveFromFolderCache(rootName string) {
	cache := service.FolderCache()

	cacheKey := fmt.Sprintf("folder:%s:%t:%t", rootName, true, false)

	cache.Delete(cacheKey)

	if err := query.UpdateAlbumFolderPreviews(); err != nil {
		log.Errorf("failed updating folder previews: %s", err)
	}

	log.Debugf("removed %s from cache", cacheKey)
}

// RemoveFromAlbumCoverCache removes covers by album UID e.g. after adding or removing photos.
func RemoveFromAlbumCoverCache(uid string) {
	cache := service.CoverCache()

	for typeName := range thumb.Types {
		cacheKey := CacheKey(albumCover, uid, typeName)

		cache.Delete(cacheKey)

		log.Debugf("removed %s from cache", cacheKey)
	}

	if err := query.UpdateAlbumPreviews(); err != nil {
		log.Errorf("failed updating album previews: %s", err)
	}
}

// FlushCoverCache clears the complete cover cache.
func FlushCoverCache() {
	service.CoverCache().Flush()

	if err := query.UpdatePreviews(); err != nil {
		log.Errorf("failed updating preview images: %s", err)
	}

	log.Debugf("albums: flushed cover cache")
}

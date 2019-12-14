package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/photoprism/photoprism/internal/util"
)

var indexer *photoprism.Indexer

func initIndexer(conf *config.Config) {
	if indexer != nil {
		return
	}

	tensorFlow := photoprism.NewTensorFlow(conf)

	indexer = photoprism.NewIndexer(conf, tensorFlow)
}

// POST /api/v1/index
func Index(router *gin.RouterGroup, conf *config.Config) {
	router.POST("/index", func(c *gin.Context) {
		if Unauthorized(c, conf) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		start := time.Now()

		var f form.IndexerOptions

		if err := c.BindJSON(&f); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.UcFirst(err.Error())})
			return
		}

		path := conf.OriginalsPath()

		event.Info(fmt.Sprintf("indexing photos in \"%s\"", filepath.Base(path)))

		if f.ConvertRaw && !conf.ReadOnly() {
			converter := photoprism.NewConverter(conf)
			converter.ConvertAll(conf.OriginalsPath())
		}

		if f.CreateThumbs {
			if err := photoprism.CreateThumbnailsFromOriginals(conf.OriginalsPath(), conf.ThumbnailsPath(), false); err != nil {
				event.Error(err.Error())
			}
		}

		initIndexer(conf)

		if f.SkipUnchanged {
			indexer.IndexOriginals(photoprism.IndexerOptionsNone())
		} else {
			indexer.IndexOriginals(photoprism.IndexerOptionsAll())
		}

		elapsed := int(time.Since(start).Seconds())

		event.Success(fmt.Sprintf("indexing completed in %d s", elapsed))
		event.Publish("index.completed", event.Data{"path": path, "seconds": elapsed})
		event.Publish("config.updated", event.Data(conf.ClientConfig()))

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("indexing completed in %d s", elapsed)})
	})
}

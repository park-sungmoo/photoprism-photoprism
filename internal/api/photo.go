package api

import (
	"fmt"
	"net/http"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/repo"
	"github.com/photoprism/photoprism/internal/util"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/photo/:uuid
//
// Parameters:
//   uuid: string PhotoUUID as returned by the API
func GetPhoto(router *gin.RouterGroup, conf *config.Config) {
	router.GET("/photos/:uuid", func(c *gin.Context) {
		if Unauthorized(c, conf) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		r := repo.New(conf.OriginalsPath(), conf.Db())
		p, err := r.PreloadPhotoByUUID(c.Param("uuid"))

		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, p)
	})
}

// PUT /api/v1/photos/:uuid
func UpdatePhoto(router *gin.RouterGroup, conf *config.Config) {
	router.PUT("/photos/:uuid", func(c *gin.Context) {
		if Unauthorized(c, conf) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		r := repo.New(conf.OriginalsPath(), conf.Db())

		m, err := r.FindPhotoByUUID(c.Param("uuid"))

		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{"error": util.UcFirst(err.Error())})
			return
		}

		if err := c.BindJSON(&m); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.UcFirst(err.Error())})
			return
		}

		conf.Db().Save(&m)

		event.Success("photo updated")

		c.JSON(http.StatusOK, m)
	})
}

// GET /api/v1/photos/:uuid/download
//
// Parameters:
//   uuid: string PhotoUUID as returned by the API
func GetPhotoDownload(router *gin.RouterGroup, conf *config.Config) {
	router.GET("/photos/:uuid/download", func(c *gin.Context) {
		r := repo.New(conf.OriginalsPath(), conf.Db())
		file, err := r.FindFileByPhotoUUID(c.Param("uuid"))

		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
			return
		}

		fileName := fmt.Sprintf("%s/%s", conf.OriginalsPath(), file.FileName)

		if !util.Exists(fileName) {
			log.Errorf("could not find original: %s", c.Param("uuid"))
			c.Data(404, "image/svg+xml", photoIconSvg)

			// Set missing flag so that the file doesn't show up in search results anymore
			file.FileMissing = true
			conf.Db().Save(&file)
			return
		}

		downloadFileName := file.DownloadFileName()

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", downloadFileName))

		c.File(fileName)
	})
}

// POST /api/v1/photos/:uuid/like
//
// Parameters:
//   uuid: string PhotoUUID as returned by the API
func LikePhoto(router *gin.RouterGroup, conf *config.Config) {
	router.POST("/photos/:uuid/like", func(c *gin.Context) {
		if Unauthorized(c, conf) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		r := repo.New(conf.OriginalsPath(), conf.Db())
		m, err := r.FindPhotoByUUID(c.Param("uuid"))

		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{"error": util.UcFirst(err.Error())})
			return
		}

		m.PhotoFavorite = true
		conf.Db().Save(&m)

		event.Publish("count.favorites", event.Data{
			"count": 1,
		})

		c.JSON(http.StatusOK, gin.H{"photo": m})
	})
}

// DELETE /api/v1/photos/:uuid/like
//
// Parameters:
//   uuid: string PhotoUUID as returned by the API
func DislikePhoto(router *gin.RouterGroup, conf *config.Config) {
	router.DELETE("/photos/:uuid/like", func(c *gin.Context) {
		if Unauthorized(c, conf) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		r := repo.New(conf.OriginalsPath(), conf.Db())
		m, err := r.FindPhotoByUUID(c.Param("uuid"))

		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{"error": util.UcFirst(err.Error())})
			return
		}

		m.PhotoFavorite = false
		conf.Db().Save(&m)

		event.Publish("count.favorites", event.Data{
			"count": -1,
		})

		c.JSON(http.StatusOK, gin.H{"photo": m})
	})
}

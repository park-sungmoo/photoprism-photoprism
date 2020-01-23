package api

import (
	"net/http"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/query"
	"github.com/photoprism/photoprism/pkg/txt"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/photoCountPerMonth
func GetPhotoCountPerMonth(router *gin.RouterGroup, conf *config.Config) {
	router.GET("/photoCountPerMonth", func(c *gin.Context) {
		if Unauthorized(c, conf) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		q := query.New(conf.OriginalsPath(), conf.Db())

		result, err := q.GetPhotoCountPerMonth()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		c.JSON(http.StatusOK, result)
	})
}

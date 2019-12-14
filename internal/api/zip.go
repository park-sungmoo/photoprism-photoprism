package api

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/repo"
	"github.com/photoprism/photoprism/internal/util"

	"github.com/gin-gonic/gin"
)

// POST /api/v1/zip
func CreateZip(router *gin.RouterGroup, conf *config.Config) {
	router.POST("/zip", func(c *gin.Context) {
		var f form.PhotoUUIDs
		start := time.Now()

		if err := c.BindJSON(&f); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.UcFirst(err.Error())})
			return
		}

		if len(f.Photos) == 0 {
			log.Error("no photos selected")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.UcFirst("no photos selected")})
			return
		}

		r := repo.New(conf.OriginalsPath(), conf.Db())
		files, err := r.FindFilesByUUID(f.Photos, 1000, 0)

		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
			return
		}

		zipPath := path.Join(conf.ExportPath(), "zip")
		zipToken := util.RandomToken(3)
		zipYear := time.Now().Format("January-2006")
		zipBaseName := fmt.Sprintf("Photos-%s-%s.zip", zipYear, zipToken)
		zipFileName := fmt.Sprintf("%s/%s", zipPath, zipBaseName)

		if err := os.MkdirAll(zipPath, 0700); err != nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": util.UcFirst("failed to create zip directory")})
			return
		}

		newZipFile, err := os.Create(zipFileName)

		if err != nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": util.UcFirst(err.Error())})
			return
		}

		defer newZipFile.Close()

		zipWriter := zip.NewWriter(newZipFile)
		defer zipWriter.Close()

		for _, file := range files {
			fileName := fmt.Sprintf("%s/%s", conf.OriginalsPath(), file.FileName)
			fileAlias := file.DownloadFileName()

			if util.Exists(fileName) {
				if err := addFileToZip(zipWriter, fileName, fileAlias); err != nil {
					log.Error(err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": util.UcFirst("failed to create zip file")})
					return
				}
				log.Infof("zip: added \"%s\" as \"%s\"", file.FileName, fileAlias)
			} else {
				log.Warnf("zip: \"%s\" is missing", file.FileName)
				file.FileMissing = true
				conf.Db().Save(&file)
			}
		}

		elapsed := int(time.Since(start).Seconds())

		log.Infof("zip: archive \"%s\" created in %s", zipBaseName, time.Since(start))

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("zip created in %d s", elapsed), "filename": zipBaseName})
	})
}

// GET /api/v1/zip/:filename
func DownloadZip(router *gin.RouterGroup, conf *config.Config) {
	router.GET("/zip/:filename", func(c *gin.Context) {
		zipBaseName := filepath.Base(c.Param("filename"))
		zipPath := path.Join(conf.ExportPath(), "zip")
		zipFileName := fmt.Sprintf("%s/%s", zipPath, zipBaseName)

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipBaseName))

		if !util.Exists(zipFileName) {
			log.Errorf("could not find zip file: %s", zipFileName)
			c.Data(404, "image/svg+xml", photoIconSvg)
			return
		}

		c.File(zipFileName)

		if err := os.Remove(zipFileName); err != nil {
			log.Errorf("zip: could not remove \"%s\" %s", zipFileName, err.Error())
		}
	})
}

func addFileToZip(zipWriter *zip.Writer, fileName, fileAlias string) error {
	fileToZip, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = fileAlias

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

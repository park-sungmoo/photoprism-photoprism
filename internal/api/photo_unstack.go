package api

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/acl"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/i18n"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/photoprism/photoprism/internal/query"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/photoprism/photoprism/pkg/txt"
)

// POST /api/v1/photos/:uid/files/:file_uid/unstack
//
// Parameters:
//   uid: string Photo UID as returned by the API
//   file_uid: string File UID as returned by the API
func PhotoUnstack(router *gin.RouterGroup) {
	router.POST("/photos/:uid/files/:file_uid/unstack", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourcePhotos, acl.ActionUpdate)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		conf := service.Config()
		fileUID := c.Param("file_uid")
		file, err := query.FileByUID(fileUID)

		if err != nil {
			log.Errorf("photo: %s (unstack)", err)
			AbortEntityNotFound(c)
			return
		}

		if file.FilePrimary {
			log.Errorf("photo: can't unstack primary file")
			AbortBadRequest(c)
			return
		} else if file.FileSidecar {
			log.Errorf("photo: can't unstack sidecar files")
			AbortBadRequest(c)
			return
		} else if file.FileRoot != entity.RootOriginals {
			log.Errorf("photo: only originals can be unstacked")
			AbortBadRequest(c)
			return
		}

		fileName := photoprism.FileName(file.FileRoot, file.FileName)
		baseName := filepath.Base(fileName)

		unstackFile, err := photoprism.NewMediaFile(fileName)

		if err != nil {
			log.Errorf("photo: %s (unstack %s)", err, txt.Quote(baseName))
			AbortEntityNotFound(c)
			return
		}

		stackPhoto := *file.Photo
		stackPrimary, err := stackPhoto.PrimaryFile()

		if err != nil {
			log.Errorf("photo: can't find primary file for existing photo (unstack %s)", txt.Quote(baseName))
			AbortUnexpected(c)
			return
		}

		related, err := unstackFile.RelatedFiles(false)

		if err != nil {
			log.Errorf("photo: %s (unstack %s)", err, txt.Quote(baseName))
			AbortEntityNotFound(c)
			return
		} else if related.Len() == 0 {
			log.Errorf("photo: no files found (unstack %s)", txt.Quote(baseName))
			AbortEntityNotFound(c)
			return
		} else if related.Main == nil {
			log.Errorf("photo: no main file found (unstack %s)", txt.Quote(baseName))
			AbortEntityNotFound(c)
			return
		}

		var files photoprism.MediaFiles
		unstackSingle := false

		if unstackFile.BasePrefix(false) == stackPhoto.PhotoName {
			if conf.ReadOnly() {
				log.Errorf("photo: can't rename files in read only mode (unstack %s)", txt.Quote(baseName))
				AbortFeatureDisabled(c)
				return
			}

			destName := fmt.Sprintf("%s.%s%s", unstackFile.AbsPrefix(false), unstackFile.Checksum(), unstackFile.Extension())

			if err := unstackFile.Move(destName); err != nil {
				log.Errorf("photo: can't rename %s to %s (unstack)", txt.Quote(unstackFile.BaseName()), txt.Quote(filepath.Base(destName)))
				AbortUnexpected(c)
				return
			}

			files = append(files, unstackFile)
			unstackSingle = true
		} else {
			files = related.Files
		}

		newPhoto := entity.NewPhoto(false)
		newPhoto.PhotoPath = unstackFile.RootRelPath()
		newPhoto.PhotoName = unstackFile.BasePrefix(false)

		if err := newPhoto.Create(); err != nil {
			log.Errorf("photo: %s (unstack %s)", err.Error(), txt.Quote(baseName))
			AbortSaveFailed(c)
			return
		}

		for _, r := range files {
			relName := r.RootRelName()
			relRoot := r.Root()

			if unstackSingle {
				relName = file.FileName
				relRoot = file.FileRoot
			}

			if err := entity.UnscopedDb().Exec(`UPDATE files 
				SET photo_id = ?, photo_uid = ?, file_name = ?, file_missing = 0
				WHERE file_name = ? AND file_root = ?`,
				newPhoto.ID, newPhoto.PhotoUID, r.RootRelName(),
				relName, relRoot).Error; err != nil {
				// Handle error...
				log.Errorf("photo: %s (unstack %s)", err.Error(), txt.Quote(r.BaseName()))

				// Remove new photo from database.
				if err := newPhoto.Delete(true); err != nil {
					log.Errorf("photo: %s (unstack %s)", err.Error(), txt.Quote(r.BaseName()))
				}

				// Revert file rename.
				if unstackSingle {
					if err := r.Move(photoprism.FileName(relRoot, relName)); err != nil {
						log.Errorf("photo: %s (unstack %s)", err.Error(), txt.Quote(r.BaseName()))
					}
				}

				AbortSaveFailed(c)
				return
			}
		}

		ind := service.Index()

		// Index unstacked files.
		if res := ind.FileName(unstackFile.FileName(), photoprism.IndexOptionsSingle()); res.Failed() {
			log.Errorf("photo: %s (unstack %s)", res.Err, txt.Quote(baseName))
			AbortSaveFailed(c)
			return
		}

		// Reset type for existing photo stack to image.
		if err := stackPhoto.Update("PhotoType", entity.TypeImage); err != nil {
			log.Errorf("photo: %s (unstack %s)", err, txt.Quote(baseName))
			AbortUnexpected(c)
			return
		}

		// Re-index existing photo stack.
		if res := ind.FileName(photoprism.FileName(stackPrimary.FileRoot, stackPrimary.FileName), photoprism.IndexOptionsSingle()); res.Failed() {
			log.Errorf("photo: %s (unstack %s)", res.Err, txt.Quote(baseName))
			AbortSaveFailed(c)
			return
		}

		// Notify clients by publishing events.
		PublishPhotoEvent(EntityCreated, newPhoto.PhotoUID, c)
		PublishPhotoEvent(EntityUpdated, stackPhoto.PhotoUID, c)

		event.SuccessMsg(i18n.MsgFileUnstacked)

		p, err := query.PhotoPreloadByUID(stackPhoto.PhotoUID)

		if err != nil {
			AbortEntityNotFound(c)
			return
		}

		c.JSON(http.StatusOK, p)
	})
}

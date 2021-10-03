package api

import (
	"fmt"
	"net/http"

	"github.com/photoprism/photoprism/internal/entity"

	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/acl"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/i18n"
	"github.com/photoprism/photoprism/internal/query"
	"github.com/photoprism/photoprism/internal/service"
)

// findFileMarker returns a file and marker entity matching the api request.
func findFileMarker(c *gin.Context) (file *entity.File, marker *entity.Marker, err error) {
	// Check authorization.
	s := Auth(SessionID(c), acl.ResourceFiles, acl.ActionUpdate)
	if s.Invalid() {
		AbortUnauthorized(c)
		return nil, nil, fmt.Errorf("unauthorized")
	}

	// Check feature flags.
	conf := service.Config()
	if !conf.Settings().Features.People {
		AbortFeatureDisabled(c)
		return nil, nil, fmt.Errorf("feature disabled")
	}

	// Find marker.
	if uid := c.Param("marker_uid"); uid == "" {
		AbortBadRequest(c)
		return nil, nil, fmt.Errorf("bad request")
	} else if marker, err = query.MarkerByUID(uid); err != nil {
		AbortEntityNotFound(c)
		return nil, nil, err
	} else if marker.FileUID == "" {
		AbortEntityNotFound(c)
		return nil, marker, fmt.Errorf("marker file missing")
	}

	// Find file.
	if f, err := query.FileByUID(marker.FileUID); err != nil {
		AbortEntityNotFound(c)
		return nil, marker, err
	} else if !f.FilePrimary {
		AbortBadRequest(c)
		return nil, marker, fmt.Errorf("can't update markers for non-primary files")
	} else {
		file = &f
	}

	return file, marker, nil
}

// UpdateMarker updates an existing file marker e.g. representing a face.
//
// PUT /api/v1/markers/:marker_uid
//
// Parameters:
//   uid: string Photo UID as returned by the API
//   file_uid: string File UID as returned by the API
//   id: int Marker ID as returned by the API
func UpdateMarker(router *gin.RouterGroup) {
	router.PUT("/markers/:marker_uid", func(c *gin.Context) {
		file, marker, err := findFileMarker(c)

		if err != nil {
			log.Debugf("marker: %s (find)", err)
			return
		}

		markerForm, err := form.NewMarker(*marker)

		if err != nil {
			log.Errorf("marker: %s (new form)", err)
			AbortSaveFailed(c)
			return
		}

		if err := c.BindJSON(&markerForm); err != nil {
			log.Errorf("marker: %s (update form)", err)
			AbortBadRequest(c)
			return
		}

		// Save marker.
		if changed, err := marker.SaveForm(markerForm); err != nil {
			log.Errorf("marker: %s", err)
			AbortSaveFailed(c)
			return
		} else if changed {
			if marker.FaceID != "" && marker.SubjUID != "" && marker.SubjSrc == entity.SrcManual {
				if res, err := service.Faces().Optimize(); err != nil {
					log.Errorf("faces: %s (optimize)", err)
				} else if res.Merged > 0 {
					log.Infof("faces: %d clusters merged", res.Merged)
				}
			}

			if err := query.UpdateSubjectCovers(); err != nil {
				log.Errorf("faces: %s (update covers)", err)
			}

			if err := entity.UpdateSubjectCounts(); err != nil {
				log.Errorf("faces: %s (update counts)", err)
			}
		}

		// Update photo metadata.
		if p, err := query.PhotoByUID(file.PhotoUID); err != nil {
			log.Errorf("faces: %s (find photo))", err)
		} else if err := p.UpdateAndSaveTitle(); err != nil {
			log.Errorf("faces: %s (update photo title)", err)
		} else {
			// Notify clients.
			PublishPhotoEvent(EntityUpdated, file.PhotoUID, c)
		}

		event.SuccessMsg(i18n.MsgChangesSaved)

		c.JSON(http.StatusOK, marker)
	})
}

// ClearMarkerSubject removes an existing marker subject association.
//
// DELETE /api/v1/markers/:marker_uid/subject
//
// Parameters:
//   uid: string Photo UID as returned by the API
//   file_uid: string File UID as returned by the API
//   id: int Marker ID as returned by the API
func ClearMarkerSubject(router *gin.RouterGroup) {
	router.DELETE("/markers/:marker_uid/subject", func(c *gin.Context) {
		file, marker, err := findFileMarker(c)

		if err != nil {
			log.Debugf("api: %s (clear marker subject)", err)
			return
		}

		if err := marker.ClearSubject(entity.SrcManual); err != nil {
			log.Errorf("faces: %s (clear subject)", err)
			AbortSaveFailed(c)
			return
		} else if err := query.UpdateSubjectCovers(); err != nil {
			log.Errorf("faces: %s (update covers)", err)
		} else if err := entity.UpdateSubjectCounts(); err != nil {
			log.Errorf("faces: %s (update counts)", err)
		}

		// Update photo metadata.
		if p, err := query.PhotoByUID(file.PhotoUID); err != nil {
			log.Errorf("faces: %s (find photo))", err)
		} else if err := p.UpdateAndSaveTitle(); err != nil {
			log.Errorf("faces: %s (update photo title)", err)
		} else {
			// Notify clients.
			PublishPhotoEvent(EntityUpdated, file.PhotoUID, c)
		}

		event.SuccessMsg(i18n.MsgChangesSaved)

		c.JSON(http.StatusOK, marker)
	})
}

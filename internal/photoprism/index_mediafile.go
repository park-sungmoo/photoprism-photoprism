package photoprism

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/photoprism/photoprism/internal/classify"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/meta"
	"github.com/photoprism/photoprism/internal/query"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/txt"
)

// MediaFile indexes a single media file.
func (ind *Index) MediaFile(m *MediaFile, o IndexOptions, originalName string) (result IndexResult) {
	if m == nil {
		err := errors.New("index: media file is nil - you might have found a bug")
		log.Error(err)
		result.Err = err
		result.Status = IndexFailed
		return result
	}

	// Skip file?
	if ind.files.Ignore(m.RootRelName(), m.Root(), m.ModTime(), o.Rescan) {
		// Skip known file.
		result.Status = IndexSkipped
		return result
	} else if o.FacesOnly && !m.IsJpeg() {
		// Skip non-jpeg file when indexing faces only.
		result.Status = IndexSkipped
		return result
	}

	start := time.Now()

	var photoQuery, fileQuery *gorm.DB
	var locKeywords []string

	file, primaryFile := entity.File{}, entity.File{}

	photo := entity.NewPhoto(o.Stack)
	metaData := meta.NewData()
	labels := classify.Labels{}
	stripSequence := Config().Settings().StackSequences() && o.Stack

	fileRoot, fileBase, filePath, fileName := m.PathNameInfo(stripSequence)
	fullBase := m.BasePrefix(false)
	logName := txt.Quote(fileName)
	fileSize, modTime, err := m.Stat()

	if err != nil {
		err := fmt.Errorf("index: %s not found", logName)
		log.Error(err)
		result.Err = err
		result.Status = IndexFailed
		return result
	}

	fileHash := ""
	fileChanged := true
	fileRenamed := false
	fileExists := false
	fileStacked := false

	photoExists := false

	event.Publish("index.indexing", event.Data{
		"fileHash": fileHash,
		"fileSize": fileSize,
		"fileName": fileName,
		"fileRoot": fileRoot,
		"baseName": filepath.Base(fileName),
	})

	// Try to find existing file by path and name.
	fileQuery = entity.UnscopedDb().First(&file, "file_name = ? AND (file_root = ? OR file_root = '')", fileName, fileRoot)
	fileExists = fileQuery.Error == nil

	// Try to find existing file by hash. Skip this for sidecar files, and files outside the originals folder.
	if !fileExists && !m.IsSidecar() && m.Root() == entity.RootOriginals {
		fileHash = m.Hash()
		fileQuery = entity.UnscopedDb().First(&file, "file_hash = ?", fileHash)

		indFileName := ""

		if fileQuery.Error == nil {
			fileExists = true
			indFileName = FileName(file.FileRoot, file.FileName)
		}

		if !fileExists {
			// Do nothing.
		} else if fs.FileExists(indFileName) {
			if err := entity.AddDuplicate(m.RootRelName(), m.Root(), m.Hash(), m.FileSize(), m.ModTime().Unix()); err != nil {
				log.Error(err)
			}

			result.Status = IndexDuplicate

			return result
		} else if err := file.Rename(m.RootRelName(), m.Root(), filePath, fileBase); err != nil {
			log.Errorf("index: %s in %s (rename)", err.Error(), logName)

			result.Status = IndexFailed
			result.Err = err

			return result
		} else if renamedSidecars, err := m.RenameSidecars(indFileName); err != nil {
			log.Errorf("index: %s in %s (rename sidecars)", err.Error(), logName)

			fileRenamed = true
		} else {
			for srcName, destName := range renamedSidecars {
				if err := query.RenameFile(entity.RootSidecar, srcName, entity.RootSidecar, destName); err != nil {
					log.Errorf("index: %s in %s (update sidecar index)", err.Error(), filepath.Join(entity.RootSidecar, srcName))
				}
			}

			fileRenamed = true
		}
	}

	// Look for existing photo if file wasn't indexed yet...
	if !fileExists {
		if photoQuery = entity.UnscopedDb().First(&photo, "photo_path = ? AND photo_name = ?", filePath, fullBase); photoQuery.Error == nil || fileBase == fullBase || !o.Stack {
			// Skip next query.
		} else if photoQuery = entity.UnscopedDb().First(&photo, "photo_path = ? AND photo_name = ? AND photo_stack > -1", filePath, fileBase); photoQuery.Error == nil {
			fileStacked = true
		}

		// Stack file based on matching location and time metadata?
		if o.Stack && photoQuery.Error != nil && Config().Settings().StackMeta() && m.MetaData().HasTimeAndPlace() {
			metaData = m.MetaData()
			photoQuery = entity.UnscopedDb().First(&photo, "photo_lat = ? AND photo_lng = ? AND taken_at = ? AND taken_src = 'meta' AND camera_serial = ?", metaData.Lat, metaData.Lng, metaData.TakenAt, metaData.CameraSerial)

			if photoQuery.Error == nil {
				fileStacked = true
			}
		}

		// Stack file based on the same unique ID?
		if o.Stack && photoQuery.Error != nil && Config().Settings().StackUUID() && m.MetaData().HasDocumentID() {
			photoQuery = entity.UnscopedDb().First(&photo, "uuid <> '' AND uuid = ?", m.MetaData().DocumentID)

			if photoQuery.Error == nil {
				fileStacked = true
			}
		}
	} else {
		photoQuery = entity.UnscopedDb().First(&photo, "id = ?", file.PhotoID)

		if fileRenamed {
			fileChanged = true
			log.Debugf("index: %s was renamed", txt.Quote(m.BaseName()))
		} else if file.Changed(fileSize, modTime) {
			fileChanged = true
			log.Debugf("index: %s was modified (new size %d, old size %d, new timestamp %d, old timestamp %d)", txt.Quote(m.BaseName()), fileSize, file.FileSize, modTime.Unix(), file.ModTime)
		} else if file.Missing() {
			fileChanged = true
			log.Debugf("index: %s was missing", txt.Quote(m.BaseName()))
		}
	}

	photoExists = photoQuery.Error == nil

	if !fileChanged && photoExists && o.SkipUnchanged() || !photoExists && m.IsSidecar() {
		result.Status = IndexSkipped
		return result
	}

	// Remove file from duplicates table if exists.
	if err := entity.PurgeDuplicate(m.RootRelName(), m.Root()); err != nil {
		log.Error(err)
	}

	details := photo.GetDetails()

	// Try to recover photo metadata from backup if not exists.
	if !photoExists {
		photo.PhotoQuality = -1

		if o.Stack {
			photo.PhotoStack = entity.IsStackable
		}

		if yamlName := fs.FormatYaml.FindFirst(m.FileName(), []string{Config().SidecarPath(), fs.HiddenPath}, Config().OriginalsPath(), stripSequence); yamlName != "" {
			if err := photo.LoadFromYaml(yamlName); err != nil {
				log.Errorf("index: %s in %s (restore from yaml)", err.Error(), logName)
			} else if err := photo.Find(); err != nil {
				log.Infof("index: %s restored from %s", txt.Quote(m.BaseName()), txt.Quote(filepath.Base(yamlName)))
			} else {
				photoExists = true
				log.Infof("index: uid %s restored from %s", photo.PhotoUID, txt.Quote(filepath.Base(yamlName)))
			}
		}
	}

	// Calculate SHA1 file hash if not exists.
	if fileHash == "" {
		fileHash = m.Hash()
	}

	// Update file hash references?
	if !fileExists || file.FileHash == "" || file.FileHash == fileHash {
		// Do nothing.
	} else if err := file.ReplaceHash(fileHash); err != nil {
		log.Errorf("index: %s while updating covers of %s", err, logName)
	}

	photo.PhotoPath = filePath

	if !o.Stack || !stripSequence || photo.PhotoStack == entity.IsUnstacked {
		photo.PhotoName = fullBase
	} else {
		photo.PhotoName = fileBase
	}

	file.FileError = ""

	// Flag first JPEG as primary file for this photo.
	if !file.FilePrimary {
		if photoExists {
			if res := entity.UnscopedDb().Where("photo_id = ? AND file_primary = 1 AND file_type = 'jpg' AND file_error = ''", photo.ID).First(&primaryFile); res.Error != nil {
				file.FilePrimary = m.IsJpeg()
			}
		} else {
			file.FilePrimary = m.IsJpeg()
		}
	}

	// Set basic file information.
	file.FileRoot = fileRoot
	file.FileName = fileName
	file.FileHash = fileHash
	file.FileSize = fileSize

	// Set file original name if available.
	if originalName != "" {
		file.OriginalName = originalName
	}

	// Set photo original name based on file original name.
	if file.OriginalName != "" {
		photo.OriginalName = fs.StripKnownExt(file.OriginalName)
	}

	if photo.PhotoQuality == -1 && (file.FilePrimary || fileChanged) {
		// Restore photos that have been purged automatically.
		photo.DeletedAt = nil
	}

	// Extra labels to ba added when new files have a photo id.
	extraLabels := classify.Labels{}

	// Detect faces in images?
	if o.FacesOnly && (!photoExists || !fileExists || !file.FilePrimary || file.FileError != "") {
		// New and non-primary files can be skipped when updating faces only.
		result.Status = IndexSkipped
		return result
	} else if ind.findFaces && file.FilePrimary {
		if markers := file.Markers(); markers != nil {
			// Detect faces.
			faces := ind.Faces(m, markers.DetectedFaceCount())

			// Create markers from faces and add them.
			if len(faces) > 0 {
				file.AddFaces(faces)
			}

			// Any new markers?
			if file.UnsavedMarkers() {
				// Add matching labels.
				extraLabels = append(extraLabels, file.Markers().Labels()...)
			} else if o.FacesOnly {
				// Skip when indexing faces only.
				result.Status = IndexSkipped
				return result
			}

			// Update photo face count.
			photo.PhotoFaces = markers.ValidFaceCount()
		} else {
			log.Errorf("index: failed loading markers for %s", logName)
		}
	}

	// Handle file types.
	switch {
	case m.IsJpeg():
		// Color information
		if p, err := m.Colors(Config().ThumbPath()); err != nil {
			log.Debugf("%s while detecting colors", err.Error())
			file.FileError = err.Error()
			file.FilePrimary = false
		} else {
			file.FileMainColor = p.MainColor.Name()
			file.FileColors = p.Colors.Hex()
			file.FileLuminance = p.Luminance.Hex()
			file.FileDiff = p.Luminance.Diff()
			file.FileChroma = p.Chroma.Value()

			if file.FilePrimary {
				photo.PhotoColor = p.MainColor.Uint8()
			}
		}

		if m.Width() > 0 && m.Height() > 0 {
			file.FileWidth = m.Width()
			file.FileHeight = m.Height()
			file.FileAspectRatio = m.AspectRatio()
			file.FilePortrait = m.Portrait()

			megapixels := m.Megapixels()

			if megapixels > photo.PhotoResolution {
				photo.PhotoResolution = megapixels
			}
		}

		if metaData := m.MetaData(); metaData.Error == nil {
			file.FileCodec = metaData.Codec
			file.SetProjection(metaData.Projection)

			if metaData.HasInstanceID() {
				log.Infof("index: %s has instance_id %s", logName, txt.Quote(metaData.InstanceID))

				file.InstanceID = metaData.InstanceID
			}
		}
	case m.IsXMP():
		if metaData, err := meta.XMP(m.FileName()); err == nil {
			// Update basic metadata.
			photo.SetTitle(metaData.Title, entity.SrcXmp)
			photo.SetDescription(metaData.Description, entity.SrcXmp)
			photo.SetTakenAt(metaData.TakenAt, metaData.TakenAtLocal, metaData.TimeZone, entity.SrcXmp)
			photo.SetCoordinates(metaData.Lat, metaData.Lng, metaData.Altitude, entity.SrcXmp)

			// Update metadata details.
			details.SetKeywords(metaData.Keywords.String(), entity.SrcXmp)
			details.SetNotes(metaData.Notes, entity.SrcXmp)
			details.SetSubject(metaData.Subject, entity.SrcXmp)
			details.SetArtist(metaData.Artist, entity.SrcXmp)
			details.SetCopyright(metaData.Copyright, entity.SrcXmp)
		} else {
			file.FileError = err.Error()
		}
	case m.IsRaw(), m.IsHEIF(), m.IsImageOther():
		if metaData := m.MetaData(); metaData.Error == nil {
			// Update basic metadata.
			photo.SetTitle(metaData.Title, entity.SrcMeta)
			photo.SetDescription(metaData.Description, entity.SrcMeta)
			photo.SetTakenAt(metaData.TakenAt, metaData.TakenAtLocal, metaData.TimeZone, entity.SrcMeta)
			photo.SetCoordinates(metaData.Lat, metaData.Lng, metaData.Altitude, entity.SrcMeta)
			photo.SetCameraSerial(metaData.CameraSerial)

			// Update metadata details.
			details.SetKeywords(metaData.Keywords.String(), entity.SrcMeta)
			details.SetNotes(metaData.Notes, entity.SrcMeta)
			details.SetSubject(metaData.Subject, entity.SrcMeta)
			details.SetArtist(metaData.Artist, entity.SrcMeta)
			details.SetCopyright(metaData.Copyright, entity.SrcMeta)

			if metaData.HasDocumentID() && photo.UUID == "" {
				log.Infof("index: %s has document_id %s", logName, txt.Quote(metaData.DocumentID))

				photo.UUID = metaData.DocumentID
			}

			if metaData.HasInstanceID() {
				log.Infof("index: %s has instance_id %s", logName, txt.Quote(metaData.InstanceID))

				file.InstanceID = metaData.InstanceID
			}

			file.FileCodec = metaData.Codec
			file.FileWidth = m.Width()
			file.FileHeight = m.Height()
			file.FileAspectRatio = m.AspectRatio()
			file.FilePortrait = m.Portrait()
			file.SetProjection(metaData.Projection)

			if res := m.Megapixels(); res > photo.PhotoResolution {
				photo.PhotoResolution = res
			}

			photo.SetCamera(entity.FirstOrCreateCamera(entity.NewCamera(m.CameraModel(), m.CameraMake())), entity.SrcMeta)
			photo.SetLens(entity.FirstOrCreateLens(entity.NewLens(m.LensModel(), m.LensMake())), entity.SrcMeta)
			photo.SetExposure(m.FocalLength(), m.FNumber(), m.Iso(), m.Exposure(), entity.SrcMeta)
		}

		if photo.TypeSrc == entity.SrcAuto {
			// Update photo type only if not manually modified.
			if m.IsRaw() && photo.PhotoType == entity.TypeImage {
				photo.PhotoType = entity.TypeRaw
			}
		}
	case m.IsVideo():
		if metaData := m.MetaData(); metaData.Error == nil {
			photo.SetTitle(metaData.Title, entity.SrcMeta)
			photo.SetDescription(metaData.Description, entity.SrcMeta)
			photo.SetTakenAt(metaData.TakenAt, metaData.TakenAtLocal, metaData.TimeZone, entity.SrcMeta)
			photo.SetCoordinates(metaData.Lat, metaData.Lng, metaData.Altitude, entity.SrcMeta)
			photo.SetCameraSerial(metaData.CameraSerial)

			// Update metadata details.
			details.SetKeywords(metaData.Keywords.String(), entity.SrcMeta)
			details.SetNotes(metaData.Notes, entity.SrcMeta)
			details.SetSubject(metaData.Subject, entity.SrcMeta)
			details.SetArtist(metaData.Artist, entity.SrcMeta)
			details.SetCopyright(metaData.Copyright, entity.SrcMeta)

			if metaData.HasDocumentID() && photo.UUID == "" {
				log.Infof("index: %s has document_id %s", logName, txt.Quote(metaData.DocumentID))

				photo.UUID = metaData.DocumentID
			}

			if metaData.HasInstanceID() {
				log.Infof("index: %s has instance_id %s", logName, txt.Quote(metaData.InstanceID))

				file.InstanceID = metaData.InstanceID
			}

			file.FileCodec = metaData.Codec
			file.FileWidth = m.Width()
			file.FileHeight = m.Height()
			file.FileAspectRatio = m.AspectRatio()
			file.FilePortrait = m.Portrait()
			file.FileDuration = metaData.Duration
			file.SetProjection(metaData.Projection)

			if res := m.Megapixels(); res > photo.PhotoResolution {
				photo.PhotoResolution = res
			}

			photo.SetCamera(entity.FirstOrCreateCamera(entity.NewCamera(m.CameraModel(), m.CameraMake())), entity.SrcMeta)
			photo.SetLens(entity.FirstOrCreateLens(entity.NewLens(m.LensModel(), m.LensMake())), entity.SrcMeta)
			photo.SetExposure(m.FocalLength(), m.FNumber(), m.Iso(), m.Exposure(), entity.SrcMeta)
		}

		if photo.TypeSrc == entity.SrcAuto {
			// Update photo type only if not manually modified.
			if file.FileDuration == 0 || file.FileDuration > time.Millisecond*3100 {
				photo.PhotoType = entity.TypeVideo
			} else {
				photo.PhotoType = entity.TypeLive
			}
		}

		if file.FileWidth == 0 && primaryFile.FileWidth > 0 {
			file.FileWidth = primaryFile.FileWidth
			file.FileHeight = primaryFile.FileHeight
			file.FileAspectRatio = primaryFile.FileAspectRatio
			file.FilePortrait = primaryFile.FilePortrait
		}

		if primaryFile.FileDiff > 0 {
			file.FileDiff = primaryFile.FileDiff
			file.FileMainColor = primaryFile.FileMainColor
			file.FileChroma = primaryFile.FileChroma
			file.FileLuminance = primaryFile.FileLuminance
			file.FileColors = primaryFile.FileColors
		}
	}

	// Set taken date based on file mod time or name if other metadata is missing.
	if m.IsMedia() && entity.SrcPriority[photo.TakenSrc] <= entity.SrcPriority[entity.SrcName] {
		// Try to extract time from original file name first.
		if taken := txt.Time(photo.OriginalName); !taken.IsZero() {
			photo.SetTakenAt(taken, taken, "", entity.SrcName)
		} else if taken, takenSrc := m.TakenAt(); takenSrc == entity.SrcName {
			photo.SetTakenAt(taken, taken, "", entity.SrcName)
		} else if !taken.IsZero() {
			photo.SetTakenAt(taken, taken, time.UTC.String(), takenSrc)
		}
	}

	// File obviously exists: remove deleted and missing flags.
	file.DeletedAt = nil
	file.FileMissing = false

	// Primary files are used for rendering thumbnails and image classification, plus sidecar files if they exist.
	if file.FilePrimary {
		primaryFile = file

		// Classify images with TensorFlow?
		if ind.findLabels {
			labels = ind.Labels(m)

			// Append labels from other sources such as face detection.
			if len(extraLabels) > 0 {
				labels = append(labels, extraLabels...)
			}

			if !photoExists && Config().Settings().Features.Private && Config().DetectNSFW() {
				photo.PhotoPrivate = ind.NSFW(m)
			}
		}

		// Read metadata from embedded Exif and JSON sidecar file, if exists.
		if metaData := m.MetaData(); metaData.Error == nil {
			// Update basic metadata.
			photo.SetTitle(metaData.Title, entity.SrcMeta)
			photo.SetDescription(metaData.Description, entity.SrcMeta)
			photo.SetTakenAt(metaData.TakenAt, metaData.TakenAtLocal, metaData.TimeZone, entity.SrcMeta)
			photo.SetCoordinates(metaData.Lat, metaData.Lng, metaData.Altitude, entity.SrcMeta)
			photo.SetCameraSerial(metaData.CameraSerial)

			// Update metadata details.
			details.SetKeywords(metaData.Keywords.String(), entity.SrcMeta)
			details.SetNotes(metaData.Notes, entity.SrcMeta)
			details.SetSubject(metaData.Subject, entity.SrcMeta)
			details.SetArtist(metaData.Artist, entity.SrcMeta)
			details.SetCopyright(metaData.Copyright, entity.SrcMeta)

			if metaData.HasDocumentID() && photo.UUID == "" {
				log.Debugf("index: %s has document_id %s", logName, txt.Quote(metaData.DocumentID))

				photo.UUID = metaData.DocumentID
			}
		}

		photo.SetCamera(entity.FirstOrCreateCamera(entity.NewCamera(m.CameraModel(), m.CameraMake())), entity.SrcMeta)
		photo.SetLens(entity.FirstOrCreateLens(entity.NewLens(m.LensModel(), m.LensMake())), entity.SrcMeta)
		photo.SetExposure(m.FocalLength(), m.FNumber(), m.Iso(), m.Exposure(), entity.SrcMeta)

		var locLabels classify.Labels

		locKeywords, locLabels = photo.UpdateLocation()
		labels = append(labels, locLabels...)
	}

	if photo.UnknownLocation() {
		photo.Cell = &entity.UnknownLocation
		photo.CellID = entity.UnknownLocation.ID
	}

	if photo.UnknownPlace() {
		photo.Place = &entity.UnknownPlace
		photo.PlaceID = entity.UnknownPlace.ID
	}

	photo.UpdateDateFields()

	if file.Panorama() {
		photo.PhotoPanorama = true
	}

	file.FileSidecar = m.IsSidecar()
	file.FileVideo = m.IsVideo()
	file.FileType = string(m.FileType())
	file.FileMime = m.MimeType()
	file.FileOrientation = m.Orientation()
	file.ModTime = modTime.Unix()

	if photoExists || photo.HasID() {
		if err := photo.Save(); err != nil {
			log.Errorf("index: %s in %s (update existing photo)", err, logName)
			result.Status = IndexFailed
			result.Err = err
			return result
		}
	} else {
		if err := photo.FirstOrCreate(); err != nil {
			log.Errorf("index: %s in %s (find or create photo)", err, logName)
			result.Status = IndexFailed
			result.Err = err
			return result
		}

		if photo.PhotoPrivate {
			event.Publish("count.private", event.Data{
				"count": 1,
			})
		}

		if photo.PhotoType == entity.TypeVideo {
			event.Publish("count.videos", event.Data{
				"count": 1,
			})
		} else {
			event.Publish("count.photos", event.Data{
				"count": 1,
			})
		}

		event.EntitiesCreated("photos", []entity.Photo{photo})
	}

	photo.AddLabels(labels)

	file.PhotoID = photo.ID
	result.PhotoID = photo.ID

	file.PhotoUID = photo.PhotoUID
	result.PhotoUID = photo.PhotoUID

	// Main JPEG file.
	if file.FilePrimary {
		labels := photo.ClassifyLabels()

		if err := photo.UpdateTitle(labels); err != nil {
			log.Debugf("%s in %s (update title)", err, logName)
		}

		w := txt.Words(details.Keywords)

		if !fs.IsGenerated(fileBase) {
			w = append(w, txt.FilenameKeywords(fileBase)...)
		}

		if photo.OriginalName == "" {
			// Do nothing.
		} else if fs.IsGenerated(photo.OriginalName) {
			w = append(w, txt.FilenameKeywords(filepath.Dir(photo.OriginalName))...)
		} else {
			w = append(w, txt.FilenameKeywords(photo.OriginalName)...)
		}

		w = append(w, txt.FilenameKeywords(filePath)...)
		w = append(w, locKeywords...)
		w = append(w, file.FileMainColor)
		w = append(w, labels.Keywords()...)

		details.Keywords = strings.Join(txt.UniqueWords(w), ", ")

		if details.Keywords != "" {
			log.Tracef("index: using keywords %s for %s", details.Keywords, logName)
		} else {
			log.Tracef("index: found no keywords for %s", logName)
		}

		photo.PhotoQuality = photo.QualityScore()

		if err := photo.Save(); err != nil {
			log.Errorf("index: %s in %s (update metadata)", err, logName)
			result.Status = IndexFailed
			result.Err = err
			return result
		}

		if err := photo.SyncKeywordLabels(); err != nil {
			log.Errorf("index: %s in %s (sync keywords and labels)", err, logName)
		}

		if err := photo.IndexKeywords(); err != nil {
			log.Errorf("index: %s in %s (save keywords)", err, logName)
		}

		if err := query.AlbumEntryFound(photo.PhotoUID); err != nil {
			log.Errorf("index: %s in %s (remove missing flag from album entry)", err, logName)
		}
	} else if err := photo.UpdateQuality(); err != nil {
		log.Errorf("index: %s in %s (update quality)", err, logName)
		result.Status = IndexFailed
		result.Err = err
		return result
	}

	result.Status = IndexUpdated

	if fileQuery.Error == nil {
		file.UpdatedIn = int64(time.Since(start))

		if err := file.Save(); err != nil {
			log.Errorf("index: %s in %s (update existing file)", err, logName)
			result.Status = IndexFailed
			result.Err = err
			return result
		}
	} else {
		file.CreatedIn = int64(time.Since(start))

		if err := file.Create(); err != nil {
			log.Errorf("index: %s in %s (add new file)", err, logName)
			result.Status = IndexFailed
			result.Err = err
			return result
		}

		event.Publish("count.files", event.Data{
			"count": 1,
		})

		if fileStacked {
			result.Status = IndexStacked
		} else {
			result.Status = IndexAdded
		}
	}

	if (photo.PhotoType == entity.TypeVideo || photo.PhotoType == entity.TypeLive) && file.FilePrimary {
		if err := file.UpdateVideoInfos(); err != nil {
			log.Errorf("index: %s in %s (update video infos)", err, logName)
		}
	}

	result.FileID = file.ID
	result.FileUID = file.FileUID

	downloadedAs := fileName

	if originalName != "" {
		downloadedAs = originalName
	}

	if err := query.SetDownloadFileID(downloadedAs, file.ID); err != nil {
		log.Errorf("index: %s in %s (set download id)", err, logName)
	}

	if !o.Stack || photo.PhotoStack == entity.IsUnstacked {
		// Do nothing.
	} else if original, merged, err := photo.Merge(Config().Settings().StackMeta(), Config().Settings().StackUUID()); err != nil {
		log.Errorf("index: %s in %s (merge)", err.Error(), logName)
	} else if len(merged) == 1 && original.ID == photo.ID {
		log.Infof("index: merged one existing photo with %s", logName)
	} else if len(merged) > 1 && original.ID == photo.ID {
		log.Infof("index: merged %d existing photos with %s", len(merged), logName)
	} else if len(merged) > 0 && original.ID != photo.ID {
		log.Infof("index: merged %s with existing photo id %d", logName, original.ID)
		result.Status = IndexStacked
		return result
	}

	if file.FilePrimary && Config().BackupYaml() {
		// Write YAML sidecar file (optional).
		yamlFile := photo.YamlFileName(Config().OriginalsPath(), Config().SidecarPath())

		if err := photo.SaveAsYaml(yamlFile); err != nil {
			log.Errorf("index: %s in %s (update yaml)", err.Error(), logName)
		} else {
			log.Debugf("index: updated yaml file %s", txt.Quote(filepath.Base(yamlFile)))
		}
	}

	return result
}

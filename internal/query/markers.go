package query

import (
	"fmt"
	"time"

	"github.com/photoprism/photoprism/internal/entity"
)

// MarkerByUID returns a Marker based on the UID.
func MarkerByUID(uid string) (*entity.Marker, error) {
	result := entity.Marker{}

	err := UnscopedDb().Where("marker_uid = ?", uid).First(&result).Error

	return &result, err
}

// Markers finds a list of file markers filtered by type, embeddings, and sorted by id.
func Markers(limit, offset int, markerType string, embeddings, subjects bool, matchedBefore time.Time) (result entity.Markers, err error) {
	db := Db()

	if markerType != "" {
		db = db.Where("marker_type = ?", markerType)
	}

	if embeddings {
		db = db.Where("embeddings_json <> ''")
	}

	if subjects {
		db = db.Where("subject_uid <> ''")
	}

	if !matchedBefore.IsZero() {
		db = db.Where("matched_at IS NULL OR matched_at < ?", matchedBefore)
	}

	db = db.Order("matched_at, marker_uid").Limit(limit).Offset(offset)

	err = db.Find(&result).Error

	return result, err
}

// UnmatchedFaceMarkers finds all currently unmatched face markers.
func UnmatchedFaceMarkers(limit, offset int, matchedBefore *time.Time) (result entity.Markers, err error) {
	db := Db().
		Where("marker_type = ?", entity.MarkerFace).
		Where("marker_invalid = 0").
		Where("embeddings_json <> ''")

	if matchedBefore == nil {
		db = db.Where("matched_at IS NULL")
	} else if !matchedBefore.IsZero() {
		db = db.Where("matched_at IS NULL OR matched_at < ?", matchedBefore)
	}

	db = db.Order("matched_at, marker_uid").Limit(limit).Offset(offset)

	err = db.Find(&result).Error

	return result, err
}

// FaceMarkers returns all face markers sorted by id.
func FaceMarkers(limit, offset int) (result entity.Markers, err error) {
	err = Db().
		Where("marker_type = ?", entity.MarkerFace).
		Order("marker_uid").Limit(limit).Offset(offset).
		Find(&result).Error

	return result, err
}

// Embeddings returns existing face embeddings.
func Embeddings(single, unclustered bool, size, score int) (result entity.Embeddings, err error) {
	var col []string

	stmt := Db().
		Model(&entity.Marker{}).
		Where("marker_type = ?", entity.MarkerFace).
		Where("marker_invalid = 0").
		Where("embeddings_json <> ''").
		Order("marker_uid")

	if size > 0 {
		stmt = stmt.Where("size >= ?", size)
	}

	if score > 0 {
		stmt = stmt.Where("score >= ?", score)
	}

	if unclustered {
		stmt = stmt.Where("face_id = ''")
	}

	if err := stmt.Pluck("embeddings_json", &col).Error; err != nil {
		return result, err
	}

	for _, embeddingsJson := range col {
		if embeddings := entity.UnmarshalEmbeddings(embeddingsJson); len(embeddings) > 0 {
			if single {
				// Single embedding per face detected.
				result = append(result, embeddings[0])
			} else {
				// Return all embedding otherwise.
				result = append(result, embeddings...)
			}
		}
	}

	return result, nil
}

// RemoveInvalidMarkerReferences removes face and subject references from invalid markers.
func RemoveInvalidMarkerReferences() (removed int64, err error) {
	res := Db().
		Model(&entity.Marker{}).
		Where("marker_invalid = 1 AND (subject_uid <> '' OR face_id <> '')").
		UpdateColumns(entity.Values{"subject_uid": "", "face_id": "", "face_dist": -1.0, "matched_at": nil})

	return res.RowsAffected, res.Error
}

// RemoveNonExistentMarkerFaces removes non-existent face IDs from the markers table.
func RemoveNonExistentMarkerFaces() (removed int64, err error) {

	res := Db().
		Model(&entity.Marker{}).
		Where("marker_type = ?", entity.MarkerFace).
		Where(fmt.Sprintf("face_id <> '' AND face_id NOT IN (SELECT id FROM %s)", entity.Face{}.TableName())).
		UpdateColumns(entity.Values{"face_id": "", "face_dist": -1.0, "matched_at": nil})

	return res.RowsAffected, res.Error
}

// RemoveNonExistentMarkerSubjects removes non-existent subject UIDs from the markers table.
func RemoveNonExistentMarkerSubjects() (removed int64, err error) {
	res := Db().
		Model(&entity.Marker{}).
		Where(fmt.Sprintf("subject_uid <> '' AND subject_uid NOT IN (SELECT subject_uid FROM %s)", entity.Subject{}.TableName())).
		UpdateColumns(entity.Values{"subject_uid": "", "matched_at": nil})

	return res.RowsAffected, res.Error
}

// FixMarkerReferences repairs invalid or non-existent references in the markers table.
func FixMarkerReferences() (removed int64, err error) {
	if r, err := RemoveInvalidMarkerReferences(); err != nil {
		return removed, err
	} else {
		removed += r
	}

	if r, err := RemoveNonExistentMarkerFaces(); err != nil {
		return removed, err
	} else {
		removed += r
	}

	if r, err := RemoveNonExistentMarkerSubjects(); err != nil {
		return removed, err
	} else {
		removed += r
	}

	return removed, nil
}

// MarkersWithNonExistentReferences finds markers with non-existent face or subject references.
func MarkersWithNonExistentReferences() (faces entity.Markers, subjects entity.Markers, err error) {
	// Find markers with invalid face IDs.
	if res := Db().
		Where("marker_type = ?", entity.MarkerFace).
		Where(fmt.Sprintf("face_id <> '' AND face_id NOT IN (SELECT id FROM %s)", entity.Face{}.TableName())).
		Find(&faces); res.Error != nil {
		err = res.Error
	}

	// Find markers with invalid subject UIDs.
	if res := Db().
		Where(fmt.Sprintf("subject_uid <> '' AND subject_uid NOT IN (SELECT subject_uid FROM %s)", entity.Subject{}.TableName())).
		Find(&subjects); res.Error != nil {
		err = res.Error
	}

	return faces, subjects, err
}

// MarkersWithSubjectConflict finds markers with conflicting subjects.
func MarkersWithSubjectConflict() (results entity.Markers, err error) {
	err = Db().
		Joins(fmt.Sprintf("JOIN %s f ON f.id = face_id AND f.subject_uid <> %s.subject_uid", entity.Face{}.TableName(), entity.Marker{}.TableName())).
		Order("face_id").
		Find(&results).Error

	return results, err
}

// ResetFaceMarkerMatches removes automatically added subject and face references from the markers table.
func ResetFaceMarkerMatches() (removed int64, err error) {
	res := Db().Model(&entity.Marker{}).
		Where("subject_src = ? AND marker_type = ?", entity.SrcAuto, entity.MarkerFace).
		UpdateColumns(entity.Values{"marker_name": "", "subject_uid": "", "subject_src": "", "face_id": "", "face_dist": -1.0, "matched_at": nil})

	return res.RowsAffected, res.Error
}

// CountUnmatchedFaceMarkers counts the number of unmatched face markers in the index.
func CountUnmatchedFaceMarkers() (n int) {
	q := Db().Model(&entity.Markers{}).
		Where("matched_at IS NULL AND marker_invalid = 0 AND embeddings_json <> ''").
		Where("marker_type = ?", entity.MarkerFace)

	if err := q.Count(&n).Error; err != nil {
		log.Errorf("faces: %s (count unmatched markers)", err)
	}

	return n
}

// CountMarkers counts the number of face markers in the index.
func CountMarkers(markerType string) (n int) {
	q := Db().Model(&entity.Markers{})

	if markerType != "" {
		q = q.Where("marker_type = ?", markerType)
	}

	if err := q.Count(&n).Error; err != nil {
		log.Errorf("faces: %s (count markers)", err)
	}

	return n
}

package search

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/dustin/go-humanize/english"
	"github.com/jinzhu/gorm"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/pluscode"
	"github.com/photoprism/photoprism/pkg/rnd"
	"github.com/photoprism/photoprism/pkg/s2"
	"github.com/photoprism/photoprism/pkg/txt"
)

// PhotosGeo searches for photos based on Form values and returns GeoResults ([]GeoResult).
func PhotosGeo(f form.PhotoSearchGeo) (results GeoResults, err error) {
	start := time.Now()

	if err := f.ParseQueryString(); err != nil {
		return GeoResults{}, err
	}

	s := UnscopedDb()

	// s.LogMode(true)

	s = s.Table("photos").
		Select(`photos.id, photos.photo_uid, photos.photo_type, photos.photo_lat, photos.photo_lng, 
		photos.photo_title, photos.photo_description, photos.photo_favorite, photos.taken_at, files.file_hash, files.file_width, 
		files.file_height`).
		Joins(`JOIN files ON files.photo_id = photos.id AND 
		files.file_missing = 0 AND files.file_primary AND files.deleted_at IS NULL`).
		Where("photos.deleted_at IS NULL").
		Where("photos.photo_lat <> 0")

	// Clip to reasonable size and normalize operators.
	f.Query = txt.NormalizeQuery(f.Query)

	// Set search filters based on search terms.
	if terms := txt.SearchTerms(f.Query); f.Query != "" && len(terms) == 0 {
		if f.Title == "" {
			f.Title = fmt.Sprintf("%s*", strings.Trim(f.Query, "%*"))
			f.Query = ""
		}
	} else if len(terms) > 0 {
		switch {
		case terms["faces"]:
			f.Query = strings.ReplaceAll(f.Query, "faces", "")
			f.Faces = "true"
		case terms["people"]:
			f.Query = strings.ReplaceAll(f.Query, "people", "")
			f.Faces = "true"
		case terms["videos"]:
			f.Query = strings.ReplaceAll(f.Query, "videos", "")
			f.Video = true
		case terms["favorites"]:
			f.Query = strings.ReplaceAll(f.Query, "favorites", "")
			f.Favorite = true
		}
	}

	// Filter by label, label category and keywords.
	if f.Query != "" {
		var categories []entity.Category
		var labels []entity.Label
		var labelIds []uint

		if err := Db().Where(AnySlug("custom_slug", f.Query, " ")).Find(&labels).Error; len(labels) == 0 || err != nil {
			log.Debugf("search: label %s not found, using fuzzy search", txt.QuoteLower(f.Query))

			for _, where := range LikeAnyKeyword("k.keyword", f.Query) {
				s = s.Where("photos.id IN (SELECT pk.photo_id FROM keywords k JOIN photos_keywords pk ON k.id = pk.keyword_id WHERE (?))", gorm.Expr(where))
			}
		} else {
			for _, l := range labels {
				labelIds = append(labelIds, l.ID)

				Db().Where("category_id = ?", l.ID).Find(&categories)

				log.Debugf("search: label %s includes %d categories", txt.QuoteLower(l.LabelName), len(categories))

				for _, category := range categories {
					labelIds = append(labelIds, category.LabelID)
				}
			}

			if wheres := LikeAnyKeyword("k.keyword", f.Query); len(wheres) > 0 {
				for _, where := range wheres {
					s = s.Where("photos.id IN (SELECT pk.photo_id FROM keywords k JOIN photos_keywords pk ON k.id = pk.keyword_id WHERE (?)) OR "+
						"photos.id IN (SELECT pl.photo_id FROM photos_labels pl WHERE pl.uncertainty < 100 AND pl.label_id IN (?))", gorm.Expr(where), labelIds)
				}
			} else {
				s = s.Where("photos.id IN (SELECT pl.photo_id FROM photos_labels pl WHERE pl.uncertainty < 100 AND pl.label_id IN (?))", labelIds)
			}
		}
	}

	// Search for one or more keywords?
	if f.Keywords != "" {
		for _, where := range LikeAnyKeyword("k.keyword", f.Keywords) {
			s = s.Where("photos.id IN (SELECT pk.photo_id FROM keywords k JOIN photos_keywords pk ON k.id = pk.keyword_id WHERE (?))", gorm.Expr(where))
		}
	}

	// Filter for specific face clusters? Example: PLJ7A3G4MBGZJRMVDIUCBLC46IAP4N7O
	if len(f.Face) >= 32 {
		for _, f := range strings.Split(strings.ToUpper(f.Face), txt.And) {
			s = s.Where(fmt.Sprintf("photos.id IN (SELECT photo_id FROM files f JOIN %s m ON f.file_uid = m.file_uid AND m.marker_invalid = 0 WHERE face_id IN (?))",
				entity.Marker{}.TableName()), strings.Split(f, txt.Or))
		}
	} else if txt.No(f.Face) {
		s = s.Where(fmt.Sprintf("photos.id IN (SELECT photo_id FROM files f JOIN %s m ON f.file_uid = m.file_uid AND m.marker_invalid = 0 AND m.marker_type = ? WHERE face_id IS NULL OR face_id = '')",
			entity.Marker{}.TableName()), entity.MarkerFace)
	} else if txt.Yes(f.Face) {
		s = s.Where(fmt.Sprintf("photos.id IN (SELECT photo_id FROM files f JOIN %s m ON f.file_uid = m.file_uid AND m.marker_invalid = 0 AND m.marker_type = ? WHERE face_id IS NOT NULL AND face_id <> '')",
			entity.Marker{}.TableName()), entity.MarkerFace)
	}

	// Filter for one or more subjects?
	if f.Subject != "" {
		for _, subj := range strings.Split(strings.ToLower(f.Subject), txt.And) {
			if subjects := strings.Split(subj, txt.Or); rnd.ContainsUIDs(subjects, 'j') {
				s = s.Where(fmt.Sprintf("photos.id IN (SELECT photo_id FROM files f JOIN %s m ON f.file_uid = m.file_uid AND m.marker_invalid = 0 WHERE subj_uid IN (?))",
					entity.Marker{}.TableName()), subjects)
			} else {
				s = s.Where(fmt.Sprintf("photos.id IN (SELECT photo_id FROM files f JOIN %s m ON f.file_uid = m.file_uid AND m.marker_invalid = 0 JOIN %s s ON s.subj_uid = m.subj_uid WHERE (?))",
					entity.Marker{}.TableName(), entity.Subject{}.TableName()), gorm.Expr(AnySlug("s.subj_slug", subj, txt.Or)))
			}
		}
	} else if f.Subjects != "" {
		for _, where := range LikeAllNames(Cols{"subj_name", "subj_alias"}, f.Subjects) {
			s = s.Where(fmt.Sprintf("photos.id IN (SELECT photo_id FROM files f JOIN %s m ON f.file_uid = m.file_uid AND m.marker_invalid = 0 JOIN %s s ON s.subj_uid = m.subj_uid WHERE (?))",
				entity.Marker{}.TableName(), entity.Subject{}.TableName()), gorm.Expr(where))
		}
	}

	// Filter by album?
	if rnd.IsPPID(f.Album, 'a') {
		s = s.Joins("JOIN photos_albums ON photos_albums.photo_uid = photos.photo_uid").
			Where("photos_albums.hidden = 0 AND photos_albums.album_uid = ?", f.Album)
	} else if f.Albums != "" || f.Album != "" {
		if f.Albums == "" {
			f.Albums = f.Album
		}

		for _, where := range LikeAnyWord("a.album_title", f.Albums) {
			s = s.Where("photos.photo_uid IN (SELECT pa.photo_uid FROM photos_albums pa JOIN albums a ON a.album_uid = pa.album_uid WHERE (?))", gorm.Expr(where))
		}
	}

	// Filter by camera?
	if f.Camera > 0 {
		s = s.Where("photos.camera_id = ?", f.Camera)
	}

	// Filter by camera lens?
	if f.Lens > 0 {
		s = s.Where("photos.lens_id = ?", f.Lens)
	}

	// Filter by year?
	if f.Year != "" {
		s = s.Where(AnyInt("photos.photo_year", f.Year, txt.Or, entity.UnknownYear, txt.YearMax))
	}

	// Filter by month?
	if f.Month != "" {
		s = s.Where(AnyInt("photos.photo_month", f.Month, txt.Or, entity.UnknownMonth, txt.MonthMax))
	}

	// Filter by day?
	if f.Day != "" {
		s = s.Where(AnyInt("photos.photo_day", f.Day, txt.Or, entity.UnknownDay, txt.DayMax))
	}

	// Find or exclude people if detected.
	if txt.IsUInt(f.Faces) {
		s = s.Where("photos.photo_faces >= ?", txt.Int(f.Faces))
	} else if txt.Yes(f.Faces) {
		s = s.Where("photos.photo_faces > 0")
	} else if txt.No(f.Faces) {
		s = s.Where("photos.photo_faces = 0")
	}

	if f.Color != "" {
		s = s.Where("files.file_main_color IN (?)", strings.Split(strings.ToLower(f.Color), txt.Or))
	}

	if f.Favorite {
		s = s.Where("photos.photo_favorite = 1")
	}

	if f.Country != "" {
		s = s.Where("photos.photo_country IN (?)", strings.Split(strings.ToLower(f.Country), txt.Or))
	}

	// Filter by media type.
	if f.Type != "" {
		s = s.Where("photos.photo_type IN (?)", strings.Split(strings.ToLower(f.Type), txt.Or))
	}

	if f.Video {
		s = s.Where("photos.photo_type = 'video'")
	} else if f.Photo {
		s = s.Where("photos.photo_type IN ('image','raw','live')")
	}

	if f.Path != "" {
		p := f.Path

		if strings.HasPrefix(p, "/") {
			p = p[1:]
		}

		if strings.HasSuffix(p, "/") {
			s = s.Where("photos.photo_path = ?", p[:len(p)-1])
		} else {
			where, values := OrLike("photos.photo_path", p)
			s = s.Where(where, values...)
		}
	}

	// Filter by primary file name without path and extension.
	if f.Name != "" {
		where, names := OrLike("photos.photo_name", f.Name)

		// Omit file path and known extensions.
		for i := range names {
			names[i] = fs.StripKnownExt(path.Base(names[i].(string)))
		}

		s = s.Where(where, names...)
	}

	// Filter by photo title.
	if f.Title != "" {
		where, values := OrLike("photos.photo_title", f.Title)
		s = s.Where(where, values...)
	}

	// Filter by status.
	if f.Archived {
		s = s.Where("photos.photo_quality > -1")
		s = s.Where("photos.deleted_at IS NOT NULL")
	} else {
		s = s.Where("photos.deleted_at IS NULL")

		if f.Private {
			s = s.Where("photos.photo_private = 1")
		} else if f.Public {
			s = s.Where("photos.photo_private = 0")
		}

		if f.Review {
			s = s.Where("photos.photo_quality < 3")
		} else if f.Quality != 0 && f.Private == false {
			s = s.Where("photos.photo_quality >= ?", f.Quality)
		}
	}

	if f.Favorite {
		s = s.Where("photos.photo_favorite = 1")
	}

	if f.S2 != "" {
		s2Min, s2Max := s2.PrefixedRange(f.S2, 7)
		s = s.Where("photos.cell_id BETWEEN ? AND ?", s2Min, s2Max)
	} else if f.Olc != "" {
		s2Min, s2Max := s2.PrefixedRange(pluscode.S2(f.Olc), 7)
		s = s.Where("photos.cell_id BETWEEN ? AND ?", s2Min, s2Max)
	} else {
		// Filter by approx distance to coordinates:
		if f.Lat != 0 {
			latMin := f.Lat - Radius*float32(f.Dist)
			latMax := f.Lat + Radius*float32(f.Dist)
			s = s.Where("photos.photo_lat BETWEEN ? AND ?", latMin, latMax)
		}
		if f.Lng != 0 {
			lngMin := f.Lng - Radius*float32(f.Dist)
			lngMax := f.Lng + Radius*float32(f.Dist)
			s = s.Where("photos.photo_lng BETWEEN ? AND ?", lngMin, lngMax)
		}
	}

	if !f.Before.IsZero() {
		s = s.Where("photos.taken_at <= ?", f.Before.Format("2006-01-02"))
	}

	if !f.After.IsZero() {
		s = s.Where("photos.taken_at >= ?", f.After.Format("2006-01-02"))
	}

	s = s.Order("taken_at, photos.photo_uid")

	if result := s.Scan(&results); result.Error != nil {
		return results, result.Error
	}

	log.Infof("geo: found %s for %s [%s]", english.Plural(len(results), "result", "results"), f.SerializeAll(), time.Since(start))

	return results, nil
}

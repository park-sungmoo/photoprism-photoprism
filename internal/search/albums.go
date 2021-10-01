package search

import (
	"strings"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/pkg/txt"
)

// Albums searches albums based on their name.
func Albums(f form.AlbumSearch) (results AlbumResults, err error) {
	if err := f.ParseQueryString(); err != nil {
		return results, err
	}

	// Base query.
	s := UnscopedDb().Table("albums").
		Select("albums.*, cp.photo_count,	cl.link_count").
		Joins("LEFT JOIN (SELECT album_uid, count(photo_uid) AS photo_count FROM photos_albums WHERE hidden = 0 AND missing = 0 GROUP BY album_uid) AS cp ON cp.album_uid = albums.album_uid").
		Joins("LEFT JOIN (SELECT share_uid, count(share_uid) AS link_count FROM links GROUP BY share_uid) AS cl ON cl.share_uid = albums.album_uid").
		Where("albums.album_type <> 'folder' OR albums.album_path IN (SELECT photo_path FROM photos WHERE photo_private = 0 AND photo_quality > -1 AND deleted_at IS NULL)").
		Where("albums.deleted_at IS NULL")

	// Limit result count.
	if f.Count > 0 && f.Count <= MaxResults {
		s = s.Limit(f.Count).Offset(f.Offset)
	} else {
		s = s.Limit(MaxResults).Offset(f.Offset)
	}

	// Set sort order.
	switch f.Order {
	case "slug":
		s = s.Order("albums.album_favorite DESC, album_slug ASC")
	default:
		s = s.Order("albums.album_favorite DESC, albums.album_year DESC, albums.album_month DESC, albums.album_day DESC, albums.album_title, albums.created_at DESC")
	}

	if f.ID != "" {
		s = s.Where("albums.album_uid IN (?)", strings.Split(f.ID, txt.Or))

		if result := s.Scan(&results); result.Error != nil {
			return results, result.Error
		}

		return results, nil
	}

	if f.Query != "" {
		likeString := "%" + f.Query + "%"
		s = s.Where("albums.album_title LIKE ? OR albums.album_location LIKE ?", likeString, likeString)
	}

	if f.Type != "" {
		s = s.Where("albums.album_type IN (?)", strings.Split(f.Type, txt.Or))
	}

	if f.Category != "" {
		s = s.Where("albums.album_category IN (?)", strings.Split(f.Category, txt.Or))
	}

	if f.Location != "" {
		s = s.Where("albums.album_location IN (?)", strings.Split(f.Location, txt.Or))
	}

	if f.Country != "" {
		s = s.Where("albums.album_country IN (?)", strings.Split(f.Country, txt.Or))
	}

	if f.Favorite {
		s = s.Where("albums.album_favorite = 1")
	}

	if (f.Year > 0 && f.Year <= txt.YearMax) || f.Year == entity.UnknownYear {
		s = s.Where("albums.album_year = ?", f.Year)
	}

	if (f.Month >= txt.MonthMin && f.Month <= txt.MonthMax) || f.Month == entity.UnknownMonth {
		s = s.Where("albums.album_month = ?", f.Month)
	}

	if (f.Day >= txt.DayMin && f.Month <= txt.DayMax) || f.Day == entity.UnknownDay {
		s = s.Where("albums.album_day = ?", f.Day)
	}

	if result := s.Scan(&results); result.Error != nil {
		return results, result.Error
	}

	return results, nil
}

package query

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/pkg/capture"
	"github.com/photoprism/photoprism/pkg/txt"
)

// Labels searches labels based on their name.
func Labels(f form.LabelSearch) (results []LabelResult, err error) {
	if err := f.ParseQueryString(); err != nil {
		return results, err
	}

	defer log.Debug(capture.Time(time.Now(), fmt.Sprintf("labels: search %s", form.Serialize(f, true))))

	s := UnscopedDb()

	// s.LogMode(true)

	s = s.Table("labels").
		Select(`labels.*`).
		Where("labels.deleted_at IS NULL").
		Where("labels.photo_count > 0").
		Group("labels.id")

	if f.ID != "" {
		s = s.Where("labels.label_uid = ?", f.ID)

		if result := s.Scan(&results); result.Error != nil {
			return results, result.Error
		}

		return results, nil
	}

	if f.Query != "" {
		var labelIds []uint
		var categories []entity.Category
		var label entity.Label

		slugString := slug.Make(f.Query)
		likeString := "%" + f.Query + "%"

		if result := Db().First(&label, "label_slug = ? OR custom_slug = ?", slugString, slugString); result.Error != nil {
			log.Infof("search: label %s not found", txt.Quote(f.Query))

			s = s.Where("labels.label_name LIKE ?", likeString)
		} else {
			labelIds = append(labelIds, label.ID)

			Db().Where("category_id = ?", label.ID).Find(&categories)

			for _, category := range categories {
				labelIds = append(labelIds, category.LabelID)
			}

			log.Infof("search: label %s includes %d categories", txt.Quote(label.LabelName), len(labelIds))

			s = s.Where("labels.id IN (?)", labelIds)
		}
	}

	if f.Favorite {
		s = s.Where("labels.label_favorite = 1")
	}

	if !f.All {
		s = s.Where("labels.label_priority >= 0 OR labels.label_favorite = 1")
	}

	switch f.Order {
	case "slug":
		s = s.Order("labels.label_favorite DESC, custom_slug ASC")
	default:
		s = s.Order("labels.label_favorite DESC, custom_slug ASC")
	}

	if f.Count > 0 && f.Count <= MaxResults {
		s = s.Limit(f.Count).Offset(f.Offset)
	} else {
		s = s.Limit(MaxResults).Offset(f.Offset)
	}

	if result := s.Scan(&results); result.Error != nil {
		return results, result.Error
	}

	return results, nil
}

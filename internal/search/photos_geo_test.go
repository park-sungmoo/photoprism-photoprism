package search

import (
	"testing"
	"time"

	"github.com/photoprism/photoprism/internal/entity"

	"github.com/photoprism/photoprism/internal/form"
	"github.com/stretchr/testify/assert"
)

func TestGeo(t *testing.T) {
	t.Run("UnknownFaces", func(t *testing.T) {
		query := form.NewGeoSearch("face:none")

		if result, err := PhotosGeo(query); err != nil {
			t.Fatal(err)
		} else {
			assert.Equal(t, 0, len(result))
		}
	})
	t.Run("form.keywords", func(t *testing.T) {
		query := form.NewGeoSearch("keywords:bridge")

		if result, err := PhotosGeo(query); err != nil {
			t.Fatal(err)
		} else {
			assert.GreaterOrEqual(t, len(result), 1)
		}
	})
	t.Run("form.subjects", func(t *testing.T) {
		query := form.NewGeoSearch("subjects:John")

		if result, err := PhotosGeo(query); err != nil {
			t.Fatal(err)
		} else {
			assert.GreaterOrEqual(t, len(result), 0)
		}
	})
	t.Run("find_all", func(t *testing.T) {
		query := form.NewGeoSearch("")

		if result, err := PhotosGeo(query); err != nil {
			t.Fatal(err)
		} else {
			assert.LessOrEqual(t, 4, len(result))
		}
	})

	t.Run("search for bridge", func(t *testing.T) {
		query := form.NewGeoSearch("Query:bridge Before:3006-01-02")
		result, err := PhotosGeo(query)
		t.Logf("RESULT: %+v", result)

		if err != nil {
			t.Fatal(err)
		}

		assert.LessOrEqual(t, 1, len(result))

	})

	t.Run("search for date range", func(t *testing.T) {
		query := form.NewGeoSearch("After:2014-12-02 Before:3006-01-02")
		result, err := PhotosGeo(query)

		// t.Logf("RESULT: %+v", result)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "Reunion", result[0].PhotoTitle)
	})

	t.Run("search for review true, quality 0", func(t *testing.T) {
		f := form.PhotoSearchGeo{
			Query:    "",
			Before:   time.Time{},
			After:    time.Time{},
			Favorite: true,
			Lat:      1.234,
			Lng:      4.321,
			S2:       "",
			Olc:      "",
			Dist:     0,
			Quality:  0,
			Review:   true,
		}

		result, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}
		assert.LessOrEqual(t, 1, len(result))
		assert.IsType(t, GeoResults{}, result)

		if len(result) > 0 {
			assert.Equal(t, "1000017", result[0].ID)
		}
	})

	t.Run("search for review false, quality > 0", func(t *testing.T) {
		f := form.PhotoSearchGeo{
			Query:    "",
			Before:   time.Time{},
			After:    time.Time{},
			Favorite: false,
			Lat:      0,
			Lng:      0,
			S2:       "",
			Olc:      "",
			Dist:     0,
			Quality:  3,
			Review:   false,
		}

		result, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}
		assert.LessOrEqual(t, 3, len(result))
		assert.IsType(t, GeoResults{}, result)
	})
	t.Run("search for s2", func(t *testing.T) {
		f := form.PhotoSearchGeo{
			Query:    "",
			Before:   time.Time{},
			After:    time.Time{},
			Favorite: false,
			Lat:      0,
			Lng:      0,
			S2:       "85",
			Olc:      "",
			Dist:     0,
			Quality:  0,
			Review:   false,
		}

		result, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, result)
		assert.IsType(t, GeoResults{}, result)
	})
	t.Run("search for Olc", func(t *testing.T) {
		f := form.PhotoSearchGeo{
			Query:    "",
			Before:   time.Time{},
			After:    time.Time{},
			Favorite: false,
			Lat:      0,
			Lng:      0,
			S2:       "",
			Olc:      "9",
			Dist:     0,
			Quality:  0,
			Review:   false,
		}

		result, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}
		assert.IsType(t, GeoResults{}, result)
	})
	t.Run("query for label flower", func(t *testing.T) {
		f := form.PhotoSearchGeo{
			Query: "flower",
		}

		result, err := PhotosGeo(f)
		if err != nil {
			t.Fatal(err)
		}

		assert.IsType(t, GeoResults{}, result)
	})
	t.Run("query for label landscape", func(t *testing.T) {
		f := form.PhotoSearchGeo{
			Query:    "landscape",
			Album:    "test",
			Camera:   123,
			Lens:     123,
			Year:     "2010",
			Month:    "12",
			Color:    "red",
			Country:  entity.UnknownID,
			Type:     "jpg",
			Video:    true,
			Path:     "/xxx/xxx/",
			Name:     "xxx",
			Archived: false,
			Private:  true,
		}

		result, err := PhotosGeo(f)
		if err != nil {
			t.Fatal(err)
		}

		assert.IsType(t, GeoResults{}, result)
	})
	t.Run("search with multiple parameters", func(t *testing.T) {
		f := form.PhotoSearchGeo{
			Query:    "landscape",
			Photo:    true,
			Path:     "/xxx,xxx",
			Name:     "xxx",
			Archived: false,
			Private:  false,
			Public:   true,
		}

		result, err := PhotosGeo(f)
		if err != nil {
			t.Fatal(err)
		}

		assert.IsType(t, GeoResults{}, result)
	})
	t.Run("search for archived true", func(t *testing.T) {
		f := form.PhotoSearchGeo{
			Query:    "landscape",
			Photo:    true,
			Path:     "/xxx/xxx/",
			Name:     "xxx",
			Archived: true,
		}

		result, err := PhotosGeo(f)
		if err != nil {
			t.Fatal(err)
		}

		assert.IsType(t, GeoResults{}, result)
	})
	t.Run("faces:true", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Query = "faces:true"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, len(photos), 4)
	})
	t.Run("faces:yes", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Faces = "Yes"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, len(photos), 4)
	})
	t.Run("faces:no", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Faces = "No"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, len(photos), 8)
	})
	t.Run("faces:2", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Faces = "2"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, len(photos), 1)
	})
	t.Run("day", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Day = "18"
		f.Month = "4"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, len(photos), 1)
	})
	t.Run("subject uid in query", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Query = "Actress"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, len(photos), 1)
	})
	t.Run("albums", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Albums = "2030"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, len(photos), 10)
	})
	t.Run("path or path", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Path = "1990/04" + "|" + "2015/11"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, len(photos), 3)
	})
	t.Run("name or name", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Name = "20151101_000000_51C501B5" + "|" + "Video"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, len(photos), 2)
	})
	t.Run("query: videos", func(t *testing.T) {
		var frm form.PhotoSearchGeo

		frm.Query = "videos"

		photos, err := PhotosGeo(frm)

		if err != nil {
			t.Fatal(err)
		}

		assert.LessOrEqual(t, 1, len(photos))

		for _, r := range photos {
			assert.IsType(t, GeoResult{}, r)
			assert.NotEmpty(t, r.ID)
			assert.Equal(t, "video", r.PhotoType)
		}
	})
	t.Run("query: faces", func(t *testing.T) {
		var frm form.PhotoSearchGeo

		frm.Query = "faces"

		photos, err := PhotosGeo(frm)

		if err != nil {
			t.Fatal(err)
		}
		assert.LessOrEqual(t, 1, len(photos))

		for _, r := range photos {
			assert.IsType(t, GeoResult{}, r)
			assert.NotEmpty(t, r.ID)
		}
	})
	t.Run("query: people", func(t *testing.T) {
		var frm form.PhotoSearchGeo

		frm.Query = "people"

		photos, err := PhotosGeo(frm)

		if err != nil {
			t.Fatal(err)
		}

		assert.LessOrEqual(t, 1, len(photos))

		for _, r := range photos {
			assert.IsType(t, GeoResult{}, r)
			assert.NotEmpty(t, r.ID)
		}
	})
	t.Run("query: favorites", func(t *testing.T) {
		var frm form.PhotoSearchGeo

		frm.Query = "favorites"

		photos, err := PhotosGeo(frm)

		if err != nil {
			t.Fatal(err)
		}

		assert.LessOrEqual(t, 1, len(photos))

		for _, r := range photos {
			assert.IsType(t, GeoResult{}, r)
			assert.NotEmpty(t, r.ID)
			assert.True(t, r.PhotoFavorite)
		}
	})
	t.Run("keywords:kuh|bridge > keywords:bridge&kuh", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Query = "keywords:kuh|bridge"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		f.Query = "keywords:bridge&kuh"

		photos2, err2 := PhotosGeo(f)

		if err2 != nil {
			t.Fatal(err2)
		}

		assert.Greater(t, len(photos), len(photos2))
	})
	t.Run("albums and and or search", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.Query = "albums:Holiday|Berlin"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		f.Query = "albums:Berlin&Holiday"

		photos2, err2 := PhotosGeo(f)

		if err2 != nil {
			t.Fatal(err2)
		}
		assert.Greater(t, len(photos), len(photos2))
	})
	t.Run("people and and or search", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.People = "Actor A|Actress A"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}

		f.People = "Actor A&Actress A"

		photos2, err2 := PhotosGeo(f)

		if err2 != nil {
			t.Fatal(err2)
		}

		assert.Greater(t, len(photos), len(photos2))
	})
	t.Run("people = subjects & person = subject", func(t *testing.T) {
		var f form.PhotoSearchGeo
		f.People = "Actor"

		photos, err := PhotosGeo(f)

		if err != nil {
			t.Fatal(err)
		}
		var f2 form.PhotoSearchGeo

		f2.Subjects = "Actor"

		photos2, err2 := PhotosGeo(f2)

		if err2 != nil {
			t.Fatal(err2)
		}

		assert.Equal(t, len(photos), len(photos2))

		var f3 form.PhotoSearchGeo

		f3.Person = "Actor A"

		photos3, err3 := PhotosGeo(f3)

		if err3 != nil {
			t.Fatal(err3)
		}

		var f4 form.PhotoSearchGeo
		f4.Subject = "Actor A"

		photos4, err4 := PhotosGeo(f4)

		if err4 != nil {
			t.Fatal(err4)
		}

		assert.Equal(t, len(photos3), len(photos4))
		assert.Equal(t, len(photos), len(photos4))
	})
}

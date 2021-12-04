package form

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGeoSearch(t *testing.T) {
	t.Run("subjects", func(t *testing.T) {
		form := &SearchGeo{Query: "subjects:\"Jens Mander\""}

		err := form.ParseQueryString()

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "Jens Mander", form.Subjects)
	})
	t.Run("aliases", func(t *testing.T) {
		form := &SearchGeo{Query: "people:\"Jens & Mander\" folder:Foo person:Bar"}

		err := form.ParseQueryString()

		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "", form.Folder)
		assert.Equal(t, "", form.Person)
		assert.Equal(t, "", form.People)
		assert.Equal(t, "Foo", form.Path)
		assert.Equal(t, "Bar", form.Subject)
		assert.Equal(t, "Jens & Mander", form.Subjects)
	})
	t.Run("keywords", func(t *testing.T) {
		form := &SearchGeo{Query: "keywords:\"Foo Bar\""}

		err := form.ParseQueryString()

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "Foo Bar", form.Keywords)
	})
	t.Run("valid query", func(t *testing.T) {
		form := &SearchGeo{Query: "query:\"fooBar baz\" before:2019-01-15 dist:25000 lat:33.45343166666667"}

		err := form.ParseQueryString()

		if err != nil {
			t.Fatal("err should be nil")
		}

		// log.Debugf("%+v\n", form)

		assert.Equal(t, "fooBar baz", form.Query)
		assert.Equal(t, time.Date(2019, 01, 15, 0, 0, 0, 0, time.UTC), form.Before)
		assert.Equal(t, uint(0x61a8), form.Dist)
		assert.Equal(t, float32(33.45343), form.Lat)
	})
	t.Run("valid query path empty folder not empty", func(t *testing.T) {
		form := &SearchGeo{Query: "query:\"fooBar baz\" before:2019-01-15 dist:25000 lat:33.45343166666667 folder:test"}

		err := form.ParseQueryString()

		if err != nil {
			t.Fatal("err should be nil")
		}

		// log.Debugf("%+v\n", form)

		assert.Equal(t, "fooBar baz", form.Query)
		assert.Equal(t, "test", form.Path)
		assert.Equal(t, "", form.Folder)
		assert.Equal(t, time.Date(2019, 01, 15, 0, 0, 0, 0, time.UTC), form.Before)
		assert.Equal(t, uint(0x61a8), form.Dist)
		assert.Equal(t, float32(33.45343), form.Lat)
	})
}

func TestGeoSearch_Serialize(t *testing.T) {
	form := &SearchGeo{Query: "query:\"fooBar baz\"", Favorite: true}

	assert.Equal(t, "q:\"query:fooBar baz\" favorite:true", form.Serialize())
}

func TestGeoSearch_SerializeAll(t *testing.T) {
	form := &SearchGeo{Query: "query:\"fooBar baz\"", Favorite: true}

	assert.Equal(t, "q:\"query:fooBar baz\" favorite:true", form.SerializeAll())
}

func TestNewGeoSearch(t *testing.T) {
	r := NewGeoSearch("Berlin")
	assert.IsType(t, SearchGeo{}, r)
}

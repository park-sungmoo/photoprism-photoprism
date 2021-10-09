package search

import (
	"testing"

	"github.com/photoprism/photoprism/internal/form"

	"github.com/stretchr/testify/assert"
)

func TestFaces(t *testing.T) {
	t.Run("Unknown", func(t *testing.T) {
		results, err := Faces(form.FaceSearch{Unknown: "yes", Order: "added", Markers: true})
		assert.NoError(t, err)
		t.Logf("Faces: %#v", results)
		if len(results) == 0 {
			t.Fatal("results are empty")
		} else if results[0].MarkerUID == "" {
			t.Fatal("marker uid is empty")
		}
	})
	t.Run("Search with limit", func(t *testing.T) {
		results, err := Faces(form.FaceSearch{Offset: 3, Order: "subject", Markers: true})
		assert.NoError(t, err)
		t.Logf("Faces: %#v", results)
		assert.LessOrEqual(t, 1, len(results))
	})
	t.Run("Find specific id", func(t *testing.T) {
		results, err := Faces(form.FaceSearch{ID: "PN6QO5INYTUSAATOFL43LL2ABAV5ACZK", Markers: true})
		assert.NoError(t, err)
		t.Logf("Faces: %#v", results)
		assert.LessOrEqual(t, 1, len(results))
	})
	t.Run("Exclude Unknown & Hidden", func(t *testing.T) {
		results, err := Faces(form.FaceSearch{Unknown: "no", Hidden: "yes", Order: "added", Markers: true})
		assert.NoError(t, err)
		t.Logf("Faces: %#v", results)
		assert.LessOrEqual(t, 0, len(results))
	})
}

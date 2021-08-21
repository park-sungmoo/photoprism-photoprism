package form

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMarker(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var m = struct {
			MarkerType    string
			MarkerSrc     string
			MarkerName    string
			SubjectUID    string
			SubjectSrc    string
			FaceID        string
			Score         int
			MarkerInvalid bool
		}{
			MarkerType:    "face",
			MarkerSrc:     "image",
			MarkerName:    "Foo",
			SubjectUID:    "3h59wvth837b5vyiub35",
			SubjectSrc:    "meta",
			FaceID:        "zz",
			Score:         100,
			MarkerInvalid: true,
		}

		f, err := NewMarker(m)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "face", f.MarkerType)
		assert.Equal(t, "image", f.MarkerSrc)
		assert.Equal(t, "Foo", f.MarkerName)
		assert.Equal(t, "3h59wvth837b5vyiub35", f.SubjectUID)
		assert.Equal(t, "zz", f.FaceID)
		assert.Equal(t, "meta", f.SubjectSrc)
		assert.Equal(t, 100, f.Score)
		assert.Equal(t, true, f.MarkerInvalid)
	})
}

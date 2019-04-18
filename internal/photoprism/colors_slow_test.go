// +build slow

package photoprism

import (
	"testing"

	"github.com/photoprism/photoprism/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestMediaFile_GetColors_Slow(t *testing.T) {
	conf := test.NewConfig()

	conf.InitializeTestData(t)

	if mediaFile2, err := NewMediaFile(conf.ImportPath() + "/iphone/IMG_6788.JPG"); err == nil {

		names, vibrantHex, mutedHex := mediaFile2.GetColors()

		t.Log(names, vibrantHex, mutedHex)

		assert.Equal(t, "#3d85c3", vibrantHex)
		assert.Equal(t, "#988570", mutedHex)
		assert.Equal(t, []string([]string{"black", "brown", "grey", "white"}), names);
	} else {
		t.Error(err)
	}

	if mediaFile3, err := NewMediaFile(conf.ImportPath() + "/raw/20140717_154212_1EC48F8489.jpg"); err == nil {

		names, vibrantHex, mutedHex := mediaFile3.GetColors()

		t.Log(names, vibrantHex, mutedHex)
		assert.Equal(t, []string([]string{"black", "brown", "grey", "white"}), names);
		assert.Equal(t, "#d5d437", vibrantHex)
		assert.Equal(t, "#a69f55", mutedHex)
	} else {
		t.Error(err)
	}
}

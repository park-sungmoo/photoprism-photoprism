// +build slow

package photoprism

import (
	"io/ioutil"
	"testing"

	"github.com/photoprism/photoprism/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestTensorFlow_GetImageTags(t *testing.T) {
	conf := test.NewConfig()

	conf.InitializeTestData(t)

	tensorFlow := NewTensorFlow(conf.TensorFlowModelPath())

	if imageBuffer, err := ioutil.ReadFile(conf.ImportPath() + "/iphone/IMG_6788.JPG"); err != nil {
		t.Error(err)
	} else {
		result, err := tensorFlow.GetImageTags(string(imageBuffer))

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.IsType(t, []TensorFlowLabel{}, result)
		assert.Equal(t, 5, len(result))

		assert.Equal(t, "tabby", result[0].Label)
		assert.Equal(t, "tiger cat", result[1].Label)

		assert.Equal(t, float32(0.16481775), result[1].Probability)
	}
}

package thumb

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"

	"github.com/photoprism/photoprism/pkg/fs"
)

// StandardRGB configures whether colors in the Apple Display P3 color space should be converted to standard RGB.
var StandardRGB = true

// Open loads an image from disk, rotates it, and converts the color profile if necessary.
func Open(fileName string, orientation int) (result image.Image, err error) {
	// Filename missing?
	if fileName == "" {
		return result, fmt.Errorf("filename missing")
	}

	// Resolve symlinks.
	if fileName, err = fs.Resolve(fileName); err != nil {
		return result, err
	}

	// Open JPEG?
	if StandardRGB && fs.FileType(fileName) == fs.ImageJPEG {
		return OpenJpeg(fileName, orientation)
	}

	// Open file with imaging function.
	img, err := imaging.Open(fileName)

	if err != nil {
		return result, err
	}

	// Adjust orientation.
	if orientation > 1 {
		img = Rotate(img, orientation)
	}

	return img, nil
}

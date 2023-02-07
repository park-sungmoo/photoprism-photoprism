package thumb

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/mandykoh/prism/meta/autometa"

	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/colors"
)

// OpenJpeg loads a JPEG image from disk, rotates it, and converts the color profile if necessary.
func OpenJpeg(fileName string, orientation int) (result image.Image, err error) {
	if fileName == "" {
		return result, fmt.Errorf("filename missing")
	}

	logName := clean.Log(filepath.Base(fileName))

	// Open file.
	fileReader, err := os.Open(fileName)

	if err != nil {
		return result, err
	}

	defer fileReader.Close()

	// Reset file offset.
	// see https://github.com/golang/go/issues/45902#issuecomment-1007953723
	_, err = fileReader.Seek(0, 0)

	if err != nil {
		return result, fmt.Errorf("%s on seek", err)
	}

	// Read color metadata.
	md, imgStream, err := autometa.Load(fileReader)

	// Decode image.
	var img image.Image

	if err != nil {
		log.Warnf("thumb: %s in %s (read color metadata)", err, logName)
		img, err = imaging.Decode(fileReader)
	} else {
		img, err = imaging.Decode(imgStream)
	}

	if err != nil {
		return result, fmt.Errorf("%s while decoding", err)
	}

	// Read ICC profile and convert colors if possible.
	if md != nil {
		if iccProfile, err := md.ICCProfile(); err != nil || iccProfile == nil {
			// Do nothing.
			log.Tracef("thumb: %s has no color profile", logName)
		} else if profile, err := iccProfile.Description(); err == nil && profile != "" {
			log.Tracef("thumb: %s has color profile %s", logName, clean.Log(profile))
			switch {
			case colors.ProfileDisplayP3.Equal(profile):
				img = colors.ToSRGB(img, colors.ProfileDisplayP3)
			}
		}
	}

	// Adjust orientation.
	if orientation > 1 {
		img = Rotate(img, orientation)
	}

	return img, nil
}

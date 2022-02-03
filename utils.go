package orientation

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/disintegration/imaging" // go get -u github.com/disintegration/imaging
	"github.com/rwcarlsen/goexif/exif"  // go get -u github.com/rwcarlsen/goexif/exif
	"github.com/sirupsen/logrus"        // go get -u github.com/sirupsen/logrus

	"bytes"
)

// ReadImage makes a copy of image (jpg,png or gif) and applies
// all necessary operation to reverse its orientation to 1 or none
// The result is a image with corrected orientation and without
// exif data.
// will also return orientation "1" or "none"
func ReadImage(imgBody []byte, logger *logrus.Entry) (imagebody []byte, orientation string) {
	imgBodyReader := bytes.NewReader(imgBody)

	// deal with exif
	img, imgExtension, err := image.Decode(imgBodyReader)
	if err != nil {
		logger.Errorf("error when decoding image, %s", err.Error())
	}
	if imgExtension != "png" && imgExtension != "jpg" && imgExtension != "jpeg" && imgExtension != "gif" {
		logger.Infof("image type %s has no exif to check for orientation", imgExtension)
		return imgBody, "none"
	}
	//dont know why, but exif needs this "hack" to decode properly sometimes
	imgBodyStringReader := strings.NewReader(string(imgBody))
	x, err := exif.Decode(imgBodyStringReader)
	if err != nil {
		if x == nil {
			logger.Warningf("Unable to read exif data, might imply that orientation is correct and no manipulation is needed, error found: %s", err)
			return imgBody, "none"
		}
		logger.Errorf("failed reading exif data: %s", err.Error())
	}
	if x != nil {
		orient, _ := x.Get(exif.Orientation)
		if orient != nil {
			if orient.String() == "1" || orient.String() == "0" {
				logger.Infof("image already has correct orientation %s, no further exif manipulation is needed", orient)
				return imgBody, orient.String()
			}
			logger.Infof("image had orientation %s", orient.String())

			img = reverseOrientation(img, orient.String(), logger)
			switch imgExtension {
			case "png":
				buffer := new(bytes.Buffer)
				err := png.Encode(buffer, img)
				if err != nil {
					logger.Errorf("error while encoding corrected image: %s", err)
				}
				imgBody = buffer.Bytes()
				return imgBody, "none"
			case "gif":
				buffer := new(bytes.Buffer)
				err := gif.Encode(buffer, img, nil)
				if err != nil {
					logger.Errorf("error while encoding corrected image: %s", err)
				}
				imgBody = buffer.Bytes()
				return imgBody, "none"
			case "jpeg", "jpg":
				buffer := new(bytes.Buffer)
				err := jpeg.Encode(buffer, img, nil)
				if err != nil {
					logger.Errorf("error while encoding corrected image: %s", err)
				}
				imgBody = buffer.Bytes()
				return imgBody, "none"
			}
		} else {
			logger.Infof("image has no orientation data - implying 1, no further exif manipulation is needed")
			return imgBody, "none"
		}
	}
	return imgBody, "none"
}

// reverseOrientation amply`s what ever operation is necessary to transform given orientation
// to the orientation 1
func reverseOrientation(img image.Image, o string, logger *logrus.Entry) *image.NRGBA {
	switch o {
	case "2":
		return imaging.FlipV(img)
	case "3":
		return imaging.Rotate180(img)
	case "4":
		return imaging.Rotate180(imaging.FlipV(img))
	case "5":
		return imaging.Rotate270(imaging.FlipV(img))
	case "6":
		return imaging.Rotate270(img)
	case "7":
		return imaging.Rotate90(imaging.FlipV(img))
	case "8":
		return imaging.Rotate90(img)
	}
	logger.Errorf("unknown orientation: %s, when attempting to rotate, expected 2-8", o)
	return imaging.Clone(img)
}

package orientation

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging" // go get -u github.com/disintegration/imaging
	"github.com/rwcarlsen/goexif/exif"  // go get -u github.com/rwcarlsen/goexif/exif
	"github.com/sirupsen/logrus"        // go get -u github.com/sirupsen/logrus

	"bytes"
)

// ReadImage makes a copy of image (jpg,png or gif) and applies
// all necessary operation to reverse its orientation to 1
// The result is a image with corrected orientation and without
// exif data.
func ReadImage(imgBody []byte) []byte {
	imgBodyReader := bytes.NewReader(imgBody)
	// deal with exif
	var err error
	img, imgExtension, err := image.Decode(imgBodyReader)
	if imgExtension != "png" && imgExtension != "jpg" && imgExtension != "jpeg" && imgExtension != "gif" {
		fmt.Printf("image type %s has no exif to check for orientation", imgExtension)
		return imgBody
	}
	x, err := exif.Decode(imgBodyReader)
	if err != nil {
		if x == nil {
			fmt.Printf("image has no exif data, no further exif manipulation is needed")
			return imgBody
		}
		logrus.Errorf("failed reading exif data: %s", err.Error())
	}
	if x != nil {
		orient, _ := x.Get(exif.Orientation)
		if orient != nil {
			if orient == 1 {
				fmt.Printf("image already has correct orientation, no further exif manipulation is needed")
				return imgBody
			}
			logrus.Infof("image had orientation %s", orient.String())
			img = reverseOrientation(img, orient.String())
			switch imgExtension {
			case "png":
				buffer := new(bytes.Buffer)
				err := png.Encode(buffer, img)
				if err != nil {
					fmt.Printf("error while encoding corrected image: %s", err)
				}
				imgBody = buffer.Bytes()
				return imgBody
			case "gif":
				buffer := new(bytes.Buffer)
				err := gif.Encode(buffer, img, nil)
				if err != nil {
					fmt.Printf("error while encoding corrected image: %s", err)
				}
				imgBody = buffer.Bytes()
				return imgBody
			case "jpeg", "jpg":
				buffer := new(bytes.Buffer)
				err := jpeg.Encode(buffer, img, nil)
				if err != nil {
					fmt.Printf("error while encoding corrected image: %s", err)
				}
				imgBody = buffer.Bytes()
				return imgBody
			}
		} else {
			logrus.Warnf("image has no orientation data - implying 1, no further exif manipulation is needed")
			return imgBody
		}
	}
	return imgBody
}

// reverseOrientation amply`s what ever operation is necessary to transform given orientation
// to the orientation 1
func reverseOrientation(img image.Image, o string) *image.NRGBA {
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
	logrus.Errorf("unknown orientation %s, expect 2-8", o)
	return imaging.Clone(img)
}

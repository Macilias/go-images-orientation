package main

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging" // go get -u github.com/disintegration/imaging
	"github.com/rwcarlsen/goexif/exif"  // go get -u github.com/rwcarlsen/goexif/exif
	"github.com/sirupsen/logrus"        // go get -u github.com/sirupsen/logrus

	"bytes"
	"errors"
	"strings"
)

// ReadImage makes a copy of image (jpg,png or gif) and applies
// all necessary operation to reverse its orientation to 1
// The result is a image with corrected orientation and without
// exif data.
func ReadImage(imgBody []byte) *image.Image {
	var err error
	// deal with exif
	if err != nil {
		logrus.Warnf("could not open file for exif decoder: %s", fpath)
	}
	imgBodyReader := bytes.NewReader(imgBody)
	img, imgExtension, err := image.Decode(imgBodyReader)
	if imgExtension != "png" && imgExtension != "jpg" && imgExtension != "gif" {
		fmt.Printf("image type %s has no exif to check for orientation", imgExtension)
		return &imgBody
	}
	x, err := exif.Decode(imgBodyReader)
	if err != nil {
		if x == nil {
			// ignore - image exif data has been already stripped
		}
		logrus.Errorf("failed reading exif data: %s", err.Error())
	}
	if x != nil {
		orient, _ := x.Get(exif.Orientation)
		if orient != nil {
			logrus.Infof("%s had orientation %s", fpath, orient.String())
			img = reverseOrientation(img, orient.String())
		} else {
			logrus.Warnf("%s had no orientation - implying 1", fpath)
			img = reverseOrientation(img, "1")
		}
		imaging.Save(img, fpath)
	}
	return &img
}

// reverseOrientation amply`s what ever operation is necessary to transform given orientation
// to the orientation 1
func reverseOrientation(img image.Image, o string) *image.NRGBA {
	switch o {
	case "1":
		return imaging.Clone(img)
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
	logrus.Errorf("unknown orientation %s, expect 1-8", o)
	return imaging.Clone(img)
}

func GetSuffix(name string) (string, error) {
	if !strings.Contains(name, ".") {
		return name, errors.New("file names without file type suffix are not supported")
	}
	split := strings.Split(name, ".")
	return strings.ToLower(strings.TrimSpace(split[len(split)-1])), nil
}

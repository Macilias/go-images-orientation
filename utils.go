package main

import (
	"image"
	"os"

	"github.com/Sirupsen/logrus"        // go get -u github.com/Sirupsen/logrus
	"github.com/disintegration/imaging" // go get -u github.com/disintegration/imaging
	"github.com/rwcarlsen/goexif/exif"  // go get -u github.com/rwcarlsen/goexif/exif

	"errors"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"
)

// ReadImage makes a copy of image (jpg,png or gif) and applies
// all necessary operation to reverse its orientation to 1
// The result is a image with corrected orientation and without
// exif data.
func ReadImage(fpath string) *image.Image {
	var img image.Image
	var err error
	// deal with image
	ifile, err := os.Open(fpath)
	if err != nil {
		logrus.Warnf("could not open file for image transformation: %s", fpath)
		return nil
	}
	defer ifile.Close()
	filetype, err := GetSuffix(fpath)
	if err != nil {
		return nil
	}
	if filetype == "jpg" {
		img, err = jpeg.Decode(ifile)
		if err != nil {
			return nil
		}
	} else if filetype == "png" {
		img, err = png.Decode(ifile)
		if err != nil {
			return nil
		}
	} else if filetype == "gif" {
		img, err = gif.Decode(ifile)
		if err != nil {
			return nil
		}
	}
	// deal with exif
	efile, err := os.Open(fpath)
	if err != nil {
		logrus.Warnf("could not open file for exif decoder: %s", fpath)
	}
	defer efile.Close()
	x, err := exif.Decode(efile)
	if err != nil {
		if x == nil {
			// ignore - image exif data has been already stripped
		}
		logrus.Errorf("failed reading exif data in [%s]: %s", fpath, err.Error())
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

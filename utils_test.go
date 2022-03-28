package orientation

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/sirupsen/logrus"
)

//This test should output the same image stripped out of exif data, can be checked on https://www.thexifer.net/ if needed
func TestF3(t *testing.T) {
	F3Img, err := getImageFromFilePath("./TestImages/F-3.jpg")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	newimg, orientation := ReadImage(F3Img, logrus.NewEntry(logrus.StandardLogger()), "whateverID1")

	os.WriteFile("./TestImages/f3_after_tests.jpg", newimg, os.ModeDevice.Perm())
	if orientation == "none" {
		t.Logf("test F3 passed with orientation %s", orientation)
		return
	}
	t.Errorf("Expected no orientation, found: %s", orientation)

}

func TestF1(t *testing.T) {
	F1Img, err := getImageFromFilePath("./TestImages/F-1.jpg")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	newimg, orientation := ReadImage(F1Img, logrus.NewEntry(logrus.StandardLogger()), "whateverID2")

	imgBodyStringReader := strings.NewReader(string(newimg))
	x, err := exif.Decode(imgBodyStringReader)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if x != nil {
		exiforientation, _ := x.Get(exif.Orientation)
		if exiforientation.String() != "1" {
			t.Errorf("Expected to read orientation 1, found %s", exiforientation)
		}
	}
	if orientation == "1" {
		return
	}
	t.Errorf("Expected orientation 1, found: %s", orientation)
}

func TestFnone(t *testing.T) {
	FnoneImg, err := getImageFromFilePath("./TestImages/F-none.jpg")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	newimg, orientation := ReadImage(FnoneImg, logrus.NewEntry(logrus.StandardLogger()), "whateverID3")

	imgBodyStringReader := strings.NewReader(string(newimg))
	x, err := exif.Decode(imgBodyStringReader)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exiforientation, _ := x.Get(exif.Orientation)
	if exiforientation != nil {
		t.Errorf("unexpected orientation found: %v", exiforientation)
	}

	if orientation != "none" {
		t.Errorf("Expected no orientation, found: %s", orientation)
	}
	t.Logf("test Fnone passed with orientation %s", orientation)
}

func getImageFromFilePath(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)

	return file, err
}

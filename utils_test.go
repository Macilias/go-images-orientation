package orientation

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

//This test should output the same image stripped out of exif data, can be checked on https://www.thexifer.net/ if needed
func TestF3(t *testing.T) {
	F3Img, err := getImageFromFilePath("./TestImages/F-3.jpg")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	newimg := ReadImage(F3Img, logrus.NewEntry(logrus.StandardLogger()), "whateverID1")

	orientation, err := GetExifOrientation(GetExif(F3Img, logrus.NewEntry(logrus.StandardLogger()), "whateverID1"))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	neworientation, err := GetExifOrientation(GetExif(newimg, logrus.NewEntry(logrus.StandardLogger()), "whateverID1"))
	if err != nil {
		t.Errorf("error: %s", err)
	}

	os.WriteFile("./TestImages/f3_after_tests.jpg", newimg, os.ModeDevice.Perm())
	if neworientation == "none" {
		t.Logf("test F3 passed with orientation from %s to %s", orientation, neworientation)
		return
	}
	t.Errorf("Expected no orientation, found: %s", neworientation)

}

func TestF1(t *testing.T) {
	F1Img, err := getImageFromFilePath("./TestImages/F-1.jpg")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	newimg := ReadImage(F1Img, logrus.NewEntry(logrus.StandardLogger()), "whateverID2")

	orientation, err := GetExifOrientation(GetExif(F1Img, logrus.NewEntry(logrus.StandardLogger()), "whateverID2"))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	neworientation, err := GetExifOrientation(GetExif(newimg, logrus.NewEntry(logrus.StandardLogger()), "whateverID2"))
	if err != nil {
		t.Errorf("error: %s", err)
	}

	if neworientation != "1" {
		t.Errorf("Expected to read orientation 1, found %s", neworientation)
	}

	if neworientation == "1" && orientation == "1" {
		t.Log("test F1 passed with unchanged orientation 1")
		return
	}
	t.Errorf("Expected orientation 1, found: %s", orientation)
}

func TestFnone(t *testing.T) {
	FnoneImg, err := getImageFromFilePath("./TestImages/F-none.jpg")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	newimg := ReadImage(FnoneImg, logrus.NewEntry(logrus.StandardLogger()), "whateverID3")

	orientation, err := GetExifOrientation(GetExif(FnoneImg, logrus.NewEntry(logrus.StandardLogger()), "whateverID3"))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	neworientation, err := GetExifOrientation(GetExif(newimg, logrus.NewEntry(logrus.StandardLogger()), "whateverID3"))
	if err != nil {
		t.Errorf("error: %s", err)
	}

	if orientation != "none" {
		t.Errorf("Expected no orientation, found: %s", orientation)
	}
	if neworientation != "none" {
		t.Errorf("unexpected orientation found: %v", neworientation)
	}

	t.Logf("test Fnone passed with orientation %s", neworientation)
}

func TestOrientError(t *testing.T) {
	var messedUpImage image.Image = nil

	t.Logf("should return error")
	orientation := "4"
	img, err := reverseOrientation(messedUpImage, orientation, logrus.NewEntry(logrus.StandardLogger()), "I D-test go :)")

	if img != nil {
		t.Errorf("unexpected result for reversing nil img")
	}

	if err == nil {
		t.Errorf("expected error was not found")
	}
}

func getImageFromFilePath(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)

	return file, err
}

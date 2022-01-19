package orientation

import (
	"fmt"
	"os"
	"testing"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func TestF3(t *testing.T) {
	F3Img, err := getImageFromFilePath("./TestImages/F-3.jpg")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	imgToSaveBytes, orientation := ReadImage(F3Img, logrus.NewEntry(logrus.StandardLogger()))

	os.WriteFile("./TestImages/f3_after_tests.jpg", imgToSaveBytes, os.ModeDevice.Perm())
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
	_, orientation := ReadImage(F1Img, logrus.NewEntry(logrus.StandardLogger()))

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
	_, orientation := ReadImage(FnoneImg, logrus.NewEntry(logrus.StandardLogger()))

	if orientation == "none" {
		t.Logf("test Fnone passed with orientation %s", orientation)
		return
	}
	t.Errorf("Expected no orientation, found: %s", orientation)
}

func getImageFromFilePath(filePath string) ([]byte, error) {
	file, err := os.Open("file.txt")
    	if err != nil {
        	log.Fatal(err)
    	}
    	defer func() {
        	if err = file.Close(); err != nil {
            		log.Fatal(err)
        }
    	}()

  	b, err := ioutil.ReadAll(file)
	return b, err
}

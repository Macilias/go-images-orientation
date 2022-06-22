Exif rotation found in tiff images can cause errors when reading images through code if not handled properly as seen in this reference: https://sirv.com/help/articles/rotate-photos-to-be-upright/


This helper builds on top of two libraries:

github.com/disintegration/imaging and
github.com/rwcarlsen/goexif/exif

```bash
go get -u github.com/disintegration/imaging
go get -u github.com/rwcarlsen/goexif/exif
```

and provides a function:

### func ReadImage(fpath string) *image.Image

which replaces the original image with a copy of referenced image (jpg, png or gif).
The replaced copy has all necessary operation, which are needed to reverse its orientation to 1, applied.
The result is a image with corrected orientation and without exif data.

Code quality: [Sonar](https://sonarcloud.io/project/overview?id=olxbr_go-images-orientation)

# Testing

Tests can be run by running `go test` on console



This helper builds on top of two libraries:

github.com/disintegration/imaging and
github.com/rwcarlsen/goexif/exif

```bash
go get -u github.com/disintegration/imaging
go get -u github.com/rwcarlsen/goexif/exif
```

and provides a function:

# func ReadImage(fpath string) *image.Image

which replaces the original image with a copy of referenced image (jpg, png or gif).
The replaced copy has all necessary operation, which are needed to reverse its orientation to 1, applied.
The result is a image with corrected orientation and without exif data.




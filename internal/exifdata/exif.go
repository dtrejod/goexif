package exifdata

import (
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func GetExifTime(filepath string) (time.Time, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return time.Time{}, err
	}

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}, err
	}

	return x.DateTime()
}

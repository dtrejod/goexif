package exifdata

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
)

func GetExifTime(filepath string) (time.Time, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return time.Time{}, err
	}
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}, err
	}

	tag, err := getDateTag(x)
	if err != nil {
		return time.Time{}, err
	}

	return formatDateFromTag(x, tag)
}

func getDateTag(x *exif.Exif) (*tiff.Tag, error) {
	tag, err := x.Get(exif.DateTimeOriginal)
	if err == nil {
		return tag, nil
	}

	tag, err = x.Get(exif.DateTimeDigitized)
	if err == nil {
		return tag, nil
	}

	return nil, errors.New("could not parse date from tiff IFD tag")
}

// copied from upstream rwcarlsen/goexif pkg
func formatDateFromTag(x *exif.Exif, tag *tiff.Tag) (time.Time, error) {
	if tag.Format() != tiff.StringVal {
		return time.Time{}, errors.New("DateTime[Original] not in string format")
	}
	exifTimeLayout := "2006:01:02 15:04:05"
	dateStr := strings.TrimRight(string(tag.Val), "\x00")
	// TODO(bradfitz,mpl): look for timezone offset, GPS time, etc.
	timeZone := time.Local
	if tz, _ := x.TimeZone(); tz != nil {
		timeZone = tz
	}
	return time.ParseInLocation(exifTimeLayout, dateStr, timeZone)
}

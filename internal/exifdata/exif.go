package exifdata

import (
	"errors"
	"os"
	"time"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

const (
	exifDateLayout = "2006:01:02 15:04:05"
)

var (
	dateTags = []string{
		"DateTimeOriginal",
		"DateTimeDigitized",
	}
)

func GetExifTime(path string) (time.Time, error) {
	// get RoofIfd
	rootIfd, err := getRootIfd(path)
	if err != nil {
		return time.Time{}, err
	}

	exifIfd, err := exif.FindIfdFromRootIfd(rootIfd, "IFD/Exif")
	if err != nil {
		return time.Time{}, errors.New("IFD/Exif not found")
	}

	value, err := getTimeFromTag(exifIfd)
	if err != nil {
		return time.Time{}, err
	}

	// Parse string into Time
	// TODO: Parse timezone
	return time.Parse(exifDateLayout, value)
}

func getTimeFromTag(exifIfd *exif.Ifd) (value string, err error) {
	for _, tag := range dateTags {
		results, err := exifIfd.FindTagWithName(tag)
		if err != nil {
			continue
		}
		if len(results) == 1 {
			return results[0].Format()
		}
	}

	return "", errors.New("could not find known IFD/Exif date tags")
}

func getRootIfd(path string) (*exif.Ifd, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rawExif, err := exif.SearchAndExtractExifWithReader(f)
	if err != nil {
		return nil, err
	}

	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return nil, err
	}
	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return nil, err
	}

	return index.RootIfd, nil
}

package moovdata

import (
	"errors"
	"os"
	"time"

	mp4 "github.com/abema/go-mp4"
)

var (
	// mvhdBoxPath is the BoxPath to the MVHD Box that contains Creation Time metadata
	// Ref: https://developer.apple.com/documentation/quicktime-file-format/movie_header_atom/creation_time
	mvhdBoxPath mp4.BoxPath = []mp4.BoxType{mp4.BoxTypeMoov(), mp4.BoxTypeMvhd()}
)

// GetTime return the CreationTime from a MooV file
func GetTime(path string) (time.Time, error) {
	boxes, err := getMetadataBoxes(path)
	if err != nil {
		return time.Time{}, err
	}
	return getTimeFromBoxes(boxes)
}

func getTimeFromBoxes(boxes []*mp4.BoxInfoWithPayload) (time.Time, error) {
	for _, box := range boxes {
		switch t := box.Payload.(type) {
		case *mp4.Mvhd:
			return timeSince1904(t.CreationTimeV0), nil
		}
	}
	return time.Time{}, errors.New("could not find mvhd box from known mp4 boxes")
}

func getMetadataBoxes(path string) ([]*mp4.BoxInfoWithPayload, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return mp4.ExtractBoxesWithPayload(f, nil, []mp4.BoxPath{mvhdBoxPath})
}

// timeSince1904 returns the time from the provided 32-bit int. CreationTime in
// mvhd4 box is the represented as seconds since 1904.
func timeSince1904(sec uint32) time.Time {
	return time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).Add(time.Second * time.Duration(sec))
}

package moovdata

import (
	"errors"
	"fmt"
	"os"
	"time"

	mp4 "github.com/abema/go-mp4"
)

var (
	// mvhdBoxPath is the BoxPath to the MVHD Box that contains Creation Time metadata
	// Ref: https://developer.apple.com/documentation/quicktime-file-format/movie_header_atom/creation_time
	mvhdBoxPath mp4.BoxPath = []mp4.BoxType{mp4.BoxTypeMoov(), mp4.BoxTypeMvhd()}
	// keysBoxPath is the BoxPath to the Keys Box. A keys box's type should
	// be intepreted before any other fields in Meta box so we can identifiy a proper index for items in ilst
	// Ref: https://developer.apple.com/documentation/quicktime-file-format/metadata_handler_atom
	keysBoxPath mp4.BoxPath = []mp4.BoxType{mp4.BoxTypeMoov(), mp4.BoxTypeMeta(), mp4.BoxTypeKeys()}
	// ilstBoxAnyPath is the BoxPath to an AnyBoxI under ilst.
	ilstBoxAnyPath mp4.BoxPath = []mp4.BoxType{mp4.BoxTypeMoov(), mp4.BoxTypeMeta(), mp4.BoxTypeIlst(), mp4.BoxTypeAny()}
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
	var creationDateIlstIdx *int
	itemIdx := 0

	fmt.Println("boxes", boxes)
	for _, box := range boxes {
		fmt.Printf("box: %+v\n", box)
		switch t := box.Payload.(type) {
		case *mp4.Keys:
			for i, k := range t.Entries {
				if string(k.KeyNamespace) == "mdta" && string(k.KeyValue) == "com.apple.quicktime.creationdate" {
					idx := i
					creationDateIlstIdx = &idx
				}
			}
		case *mp4.Data:
			fmt.Println("item")
			itemIdx++
			if creationDateIlstIdx != nil && itemIdx == *creationDateIlstIdx {
				fmt.Println(string(t.Data))
			}
		case *mp4.Mvhd:
			// TODO: Use mvhd block creation date as fallback
			//return timeSince1904(t.CreationTimeV0), nil
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

	return mp4.ExtractBoxesWithPayload(f, nil, []mp4.BoxPath{
		keysBoxPath,
		ilstBoxAnyPath,
		mvhdBoxPath})
}

// timeSince1904 returns the time from the provided 32-bit int. CreationTime in
// mvhd4 box is the represented as seconds since 1904.
func timeSince1904(sec uint32) time.Time {
	return time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).Add(time.Second * time.Duration(sec))
}

package moovdata

import (
	"errors"
	"fmt"
	"os"
	"time"

	mp4 "github.com/abema/go-mp4"
)

const (
	ilstItemDateLayout = "2006-01-02T15:04:05-0700"
)

var (
	// handledBoxes are the set of boxes that we will expand
	handledBoxes map[mp4.BoxType]struct{} = map[mp4.BoxType]struct{}{
		mp4.BoxTypeMoov(): {},
		// mvhd box will contain creation time
		// Ref: https://developer.apple.com/documentation/quicktime-file-format/movie_header_atom/creation_time
		mp4.BoxTypeMvhd(): {},
		// meta box will contain ilst item where 1 item may contain creation time
		mp4.BoxTypeMeta(): {},
		mp4.BoxTypeKeys(): {},
		mp4.BoxTypeIlst(): {},
	}
	// handledChildBoxes are the set of boxes that we will expand if the
	// parent was one of the following from the BoxPath. Any item in this
	// list should also be in handledBoxes in order to pickup the children
	// of the box
	handledChildBoxes map[mp4.BoxType]struct{} = map[mp4.BoxType]struct{}{
		mp4.BoxTypeIlst(): {},
	}
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

	var ilstTs *time.Time
	var mvhdTs *time.Time
	for _, box := range boxes {
		switch t := box.Payload.(type) {
		case *mp4.Keys:
			for i, k := range t.Entries {
				if string(k.KeyNamespace) == "mdta" && string(k.KeyValue) == "com.apple.quicktime.creationdate" {
					// ilst items are not zero indexed so add 1
					idx := i + 1
					creationDateIlstIdx = &idx
				}
			}
		case *mp4.Item:
			itemIdx++
			if creationDateIlstIdx != nil && itemIdx == *creationDateIlstIdx {
				ts, err := time.Parse(ilstItemDateLayout, string(t.Data.Data[:]))
				if err != nil {
					// unexpected error but swallow in case we find a mvhd timestamp
					continue
				}
				ilstTs = &ts
			}
		case *mp4.Mvhd:
			ts := timeSince1904(t.CreationTimeV0)
			mvhdTs = &ts
		default:
		}
	}

	// prefer ilst timestamp
	if ilstTs != nil {
		return (*ilstTs).UTC(), nil
	}
	if mvhdTs != nil {
		fmt.Println("mvhd")
		return (*mvhdTs).UTC(), nil
	}

	return time.Time{}, errors.New("could not find box from known mp4 boxes that contains a valid timestamp")
}

func getMetadataBoxes(path string) ([]*mp4.BoxInfoWithPayload, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var handledChild bool
	var bis []*mp4.BoxInfo
	_, err = mp4.ReadBoxStructure(f, func(h *mp4.ReadHandle) (interface{}, error) {
		_, handled := handledBoxes[h.BoxInfo.Type]

		if len(h.Path) > 1 {
			_, handledChild = handledChildBoxes[h.Path[len(h.Path)-2]]
		} else {
			handledChild = false
		}

		if handled || handledChild {
			// set if underIlst on box context for parsing box code below
			if len(h.Path) > 0 {
				if h.Path[len(h.Path)-1] == mp4.BoxTypeIlst() {
					h.BoxInfo.Context.UnderIlst = true
				}
			}

			bis = append(bis, &h.BoxInfo)

			// Expands children
			return h.Expand()
		}
		return nil, nil
	})

	bs := make([]*mp4.BoxInfoWithPayload, 0, len(bis))
	for _, bi := range bis {
		if _, err := bi.SeekToPayload(f); err != nil {
			return nil, err
		}

		box, _, err := mp4.UnmarshalAny(f, bi.Type, bi.Size-bi.HeaderSize, bi.Context)
		if err != nil {
			return nil, err
		}
		bs = append(bs, &mp4.BoxInfoWithPayload{
			Info:    *bi,
			Payload: box,
		})
	}
	return bs, nil
}

// timeSince1904 returns the time from the provided 32-bit int. CreationTime in
// mvhd4 box is the represented as seconds since 1904.
func timeSince1904(sec uint32) time.Time {
	return time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).Add(time.Second * time.Duration(sec))
}

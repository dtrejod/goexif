package riffdata

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/go-errors/errors"
	"golang.org/x/image/riff"
)

const (
	// riffIDITDateLayout is the string time layout used by the IDIT chunk
	riffIDITDateLayout = time.ANSIC
)

var (
	// idit is the riff tag associated with the DateTimeOriginal riff tag
	// Ref: https://exiftool.org/TagNames/RIFF.html
	idit = riff.FourCC{'I', 'D', 'I', 'T'}
)

// GetTime returns the RIFF metadata Datetime from media referenced in the provided path
func GetTime(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()

	_, r, err := riff.NewReader(f)
	if err != nil {
		return time.Time{}, errors.WrapPrefix(err, "could not open file with RIFF reader", 0)
	}
	t, found, err := dumpTimeFromReader(r)
	if err != nil {
		return time.Time{}, errors.WrapPrefix(err, "unexpected error parsing datetime from RIFF metadata", 0)
	}
	if !found {
		return time.Time{}, errors.New("could not find datetime in RIFF metadata")
	}
	return t, nil
}

func dumpTimeFromReader(r *riff.Reader) (time.Time, bool, error) {
	for {
		chunkID, chunkLen, chunkData, err := r.Next()
		if err == io.EOF {
			return time.Time{}, false, nil
		}
		if err != nil {
			return time.Time{}, false, err
		}

		switch chunkID {
		case riff.LIST:
			_, listChunk, err := riff.NewListReader(chunkLen, chunkData)
			if err != nil {
				return time.Time{}, false, err
			}
			t, found, err := dumpTimeFromReader(listChunk)
			if err != nil {
				return time.Time{}, false, err
			}
			if !found {
				continue
			}
			return t, true, nil
		case idit:
			t, err := getTimeFromReader(chunkData, chunkLen)
			return t, true, err
		}
	}
}

func getTimeFromReader(r io.Reader, chunkLen uint32) (time.Time, error) {
	buf := bytes.NewBuffer(make([]byte, 0, chunkLen))
	_, err := io.Copy(buf, r)
	if err != nil {
		return time.Time{}, errors.WrapPrefix(err, "failed to copy RIFF IDIT tag into buffer", 0)
	}

	// IDIT chunk is both line feed and null char terminated
	val := string(bytes.TrimRight(buf.Bytes(), "\x0a\x00"))

	// Parse string into Time
	t, err := time.Parse(riffIDITDateLayout, val)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

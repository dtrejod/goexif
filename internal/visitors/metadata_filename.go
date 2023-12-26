package visitors

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dtrejod/goexif/internal/exifdata"
	"github.com/dtrejod/goexif/internal/mediatype"
)

const (
	outPathDateFormat = "2006/01/02"
)

type mediaMetadataFilename struct {
	outDir              *string
	useLastModifiedDate bool
	timestampAsFilename bool
	cleanOutFileExt     bool
}

// NewMediaMetadataFilename is a mediatype visitor that will generated an appropriate
// output filename for a provided mediatype format.
// - useLastModifiedDate: Fallback to using the last modified date if no EXIF data exists on the media
// - timestampAsFilename: Use the Unix EPOCH time as the output file name.
// - cleanOutFileExt: Use the identified mediatype Ext as the extension of the output filename
// TODO(dtrejo): Rename useLastModifiedDate to fallbackToLastModifiedDate to
// better describe what this variable actually does.
func NewMediaMetadataFilename(
	_ context.Context,
	outDir *string,
	useLastModifiedDate,
	timestampAsFilename,
	cleanOutFileExt bool,
) mediatype.VisitorFunc[string] {
	return &mediaMetadataFilename{
		outDir:              outDir,
		useLastModifiedDate: useLastModifiedDate,
		timestampAsFilename: timestampAsFilename,
		cleanOutFileExt:     cleanOutFileExt,
	}
}

func (e *mediaMetadataFilename) VisitJPEG(ctx context.Context, image mediatype.JPEG) (string, error) {
	ts, err := exifdata.GetExifTime(image.Path)
	if err != nil {
		ts, err = e.fallbackToModTime(image.Path, err)
		if err != nil {
			return "", err
		}
	}
	return e.getOutputFile(ctx, image.Path, image.Ext(), ts.UTC())
}

// VisitPNG implements VisitorFunc
// EXIF extension was adopted for PNG in 2017
// http://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html#C.eXIf
func (e *mediaMetadataFilename) VisitPNG(ctx context.Context, image mediatype.PNG) (string, error) {
	ts, err := exifdata.GetExifTime(image.Path)
	if err != nil {
		ts, err = e.fallbackToModTime(image.Path, err)
		if err != nil {
			return "", err
		}
	}
	return e.getOutputFile(ctx, image.Path, image.Ext(), ts.UTC())
}

func (e *mediaMetadataFilename) VisitHEIF(ctx context.Context, image mediatype.HEIF) (string, error) {
	return "", nil
}

func (e *mediaMetadataFilename) getOutputFile(_ context.Context, srcPath, cleanExt string, ts time.Time) (string, error) {
	srcDir := filepath.Dir(srcPath)
	outDir := filepath.Join(srcDir, ts.Format(outPathDateFormat))
	if e.outDir != nil {
		outDir = filepath.Join(*e.outDir, ts.Format(outPathDateFormat))
	}

	ext := filepath.Ext(strings.ToLower(srcPath))
	if e.cleanOutFileExt {
		ext = cleanExt
	}
	outFilename := strings.TrimSuffix(filepath.Base(srcPath), filepath.Ext(srcPath))
	if e.timestampAsFilename {
		outFilename = strconv.FormatInt(ts.Unix(), 10)
	}

	outFilename = outFilename + ext
	return filepath.Join(outDir, outFilename), nil

}

func (e *mediaMetadataFilename) fallbackToModTime(srcPath string, origErr error) (time.Time, error) {
	// on error, fallback to lastmodified if the option was specified
	if e.useLastModifiedDate {
		f, statErr := os.Stat(srcPath)
		if statErr != nil {
			return time.Time{}, origErr
		}
		return f.ModTime(), nil
	}
	return time.Time{}, origErr
}

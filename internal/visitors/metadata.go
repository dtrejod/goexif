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
	outDir                  *string
	useLastModifiedDate     bool
	timestampAsFilename     bool
	useOutputMagicSignature bool
}

// MediaMetadata is the return type from the MediaMetadataFilename visitor
type MediaMetadata struct {
	// OutPath is an appropriate new output filename for the provided mediatype format.
	OutPath   string
	Timestamp time.Time
}

// NewMediaMetadataFilename is a mediatype visitor that will generate metadata info on a provided media file
// - useLastModifiedDate: Fallback to using the last modified date if no EXIF data exists on the media
// - timestampAsFilename: Use the Unix EPOCH time as the output file name.
// - useOutputMagicSignature: Use the identified mediatype Ext as the extension of the output filename
// TODO(dtrejo): Rename useLastModifiedDate to fallbackToLastModifiedDate to
// better describe what this variable actually does.
func NewMediaMetadataFilename(
	_ context.Context,
	outDir *string,
	useLastModifiedDate,
	timestampAsFilename,
	useOutputMagicSignature bool,
) mediatype.VisitorFunc[MediaMetadata] {
	return &mediaMetadataFilename{
		outDir:                  outDir,
		useLastModifiedDate:     useLastModifiedDate,
		timestampAsFilename:     timestampAsFilename,
		useOutputMagicSignature: useOutputMagicSignature,
	}
}

func (e *mediaMetadataFilename) VisitJPEG(ctx context.Context, image mediatype.JPEG) (MediaMetadata, error) {
	return e.getEXIFMetadata(ctx, image.Path, image.Ext())
}

// VisitPNG implements VisitorFunc
// EXIF extension was adopted for PNG in 2017
// http://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html#C.eXIf
func (e *mediaMetadataFilename) VisitPNG(ctx context.Context, image mediatype.PNG) (MediaMetadata, error) {
	return e.getEXIFMetadata(ctx, image.Path, image.Ext())
}

func (e *mediaMetadataFilename) VisitHEIF(ctx context.Context, image mediatype.HEIF) (MediaMetadata, error) {
	return e.getEXIFMetadata(ctx, image.Path, image.Ext())
}

func (e *mediaMetadataFilename) VisitTIFF(ctx context.Context, image mediatype.TIFF) (MediaMetadata, error) {
	return e.getEXIFMetadata(ctx, image.Path, image.Ext())
}

func (e *mediaMetadataFilename) getEXIFMetadata(ctx context.Context, srcPath, cleanEXT string) (MediaMetadata, error) {
	ts, err := exifdata.GetExifTime(srcPath)
	if err != nil {
		ts, err = e.fallbackToModTime(srcPath, err)
		if err != nil {
			return MediaMetadata{}, err
		}
	}
	outFile, err := e.getOutputFile(ctx, srcPath, cleanEXT, ts.UTC())
	if err != nil {
		return MediaMetadata{}, err
	}

	return MediaMetadata{
		OutPath:   outFile,
		Timestamp: ts,
	}, nil
}

func (e *mediaMetadataFilename) getOutputFile(_ context.Context, srcPath, cleanExt string, ts time.Time) (string, error) {
	srcDir := filepath.Dir(srcPath)
	outDir := filepath.Join(srcDir, ts.Format(outPathDateFormat))
	if e.outDir != nil {
		outDir = filepath.Join(*e.outDir, ts.Format(outPathDateFormat))
	}

	ext := filepath.Ext(srcPath)
	if e.useOutputMagicSignature {
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

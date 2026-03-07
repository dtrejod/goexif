package visitors

import (
	"context"
	"os"
	"time"

	"github.com/dtrejod/goexif/internal/exifdata"
	"github.com/dtrejod/goexif/internal/mediatype"
	"github.com/dtrejod/goexif/internal/moovdata"
	"github.com/dtrejod/goexif/internal/riffdata"
)

const (
	outPathDateFormat = "2006/01/02"
)

type mediaMetadataFilename struct {
	useLastModifiedDate bool
	timestampAsFilename bool
}

// MediaMetadata is the return type from the MediaMetadataFilename visitor
type MediaMetadata struct {
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
		useLastModifiedDate: useLastModifiedDate,
		timestampAsFilename: timestampAsFilename,
	}
}

func (e *mediaMetadataFilename) VisitJPEG(ctx context.Context, image mediatype.JPEG) (MediaMetadata, error) {
	return e.getTimeMetadataWithFunc(ctx, image.Path, exifdata.GetTime, image.Ext())
}

// VisitPNG implements VisitorFunc
// EXIF extension was adopted for PNG in 2017
// http://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html#C.eXIf
func (e *mediaMetadataFilename) VisitPNG(ctx context.Context, image mediatype.PNG) (MediaMetadata, error) {
	return e.getTimeMetadataWithFunc(ctx, image.Path, exifdata.GetTime, image.Ext())
}

func (e *mediaMetadataFilename) VisitHEIF(ctx context.Context, image mediatype.HEIF) (MediaMetadata, error) {
	return e.getTimeMetadataWithFunc(ctx, image.Path, exifdata.GetTime, image.Ext())
}

func (e *mediaMetadataFilename) VisitTIFF(ctx context.Context, image mediatype.TIFF) (MediaMetadata, error) {
	return e.getTimeMetadataWithFunc(ctx, image.Path, exifdata.GetTime, image.Ext())
}

func (e *mediaMetadataFilename) VisitQTFF(ctx context.Context, image mediatype.QTFF) (MediaMetadata, error) {
	return e.getTimeMetadataWithFunc(ctx, image.Path, moovdata.GetTime, image.Ext())
}

func (e *mediaMetadataFilename) VisitMP4(ctx context.Context, image mediatype.MP4) (MediaMetadata, error) {
	return e.getTimeMetadataWithFunc(ctx, image.Path, moovdata.GetTime, image.Ext())
}

func (e *mediaMetadataFilename) VisitAVI(ctx context.Context, image mediatype.AVI) (MediaMetadata, error) {
	return e.getTimeMetadataWithFunc(ctx, image.Path, riffdata.GetTime, image.Ext())
}

func (e *mediaMetadataFilename) Visit3PG(ctx context.Context, image mediatype.GPP) (MediaMetadata, error) {
	return e.getTimeMetadataWithFunc(ctx, image.Path, moovdata.GetTime, image.Ext())
}

func (e *mediaMetadataFilename) Visit3G2(ctx context.Context, image mediatype.GPP2) (MediaMetadata, error) {
	return e.getTimeMetadataWithFunc(ctx, image.Path, moovdata.GetTime, image.Ext())
}

func (e *mediaMetadataFilename) getTimeMetadataWithFunc(
	ctx context.Context,
	srcPath string,
	tsFunc func(string) (time.Time, error),
	cleanEXT string,
) (MediaMetadata, error) {
	ts, err := tsFunc(srcPath)
	if err != nil {
		ts, err = e.fallbackToModTime(srcPath, err)
		if err != nil {
			return MediaMetadata{}, err
		}
	}

	return MediaMetadata{
		Timestamp: ts,
	}, nil
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

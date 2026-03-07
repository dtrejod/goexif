package visitors

import "github.com/dtrejod/goexif/internal/mediatype"

type renameMedia struct {
	outDir                  *string
	useOutputMagicSignature bool

	mediaDaterVisitor mediatype.VisitorFunc[MediaMetadata]
}

// NewRenameMedia is a visitor that will rename a media file based on the metadata
// contained within the file.
func NewRenameMedia(
	outDir *string,
	useOutputMagicSignature bool,
	mediaDaterVisitor mediatype.VisitorFunc[MediaMetadata],
) mediatype.VisitorFunc[MediaMetadata] {
	return &renameMedia{
		outDir:                  outDir,
		useOutputMagicSignature: useOutputMagicSignature,
		mediaDaterVisitor:       mediaDaterVisitor,
	}
}

func (e *renameMedia) VisitJPEG(ctx context.Context, image mediatype.JPEG) (MediaMetadata, error) {
	return e.renameMedia(ctx, image.Path, image.Ext())
}

func (e *renameMedia) VisitPNG(ctx context.Context, image mediatype.PNG) (MediaMetadata, error) {
	return e.renameMedia(ctx, image.Path, image.Ext())
}


func (e *renameMedia) getOutputFileName(
	ctx context.Context,
	srcPath string,
	cleanEXT string,
) (string, error) {
	// Get the metadata from the media file

	return.getOutputFile(ctx, srcPath, cleanEXT, metadata.TimeStamp.UTC())


func (e *renameMedia) getOutputFile(_ context.Context, srcPath, cleanExt string, ts time.Time) (string, error) {
	metadata, err := e.mediaDaterVisitor(ctx, srcPath)
	if err != nil {
		return "", err
	}

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


package mediasort

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dtrejod/goexif/internal/ilog"
	"github.com/dtrejod/goexif/internal/mediatype"
	"github.com/dtrejod/goexif/internal/visitors"
	"go.uber.org/zap"
)

type metadataFileHandler struct {
	useInputMagicSignature bool
	dryRun                 bool
	detectDuplicates       bool
	overwriteExisting      bool

	mediaMetadataVisitorFunc mediatype.VisitorFunc[visitors.MediaMetadata]
}

// handle takes in a source media file, and will move it a computed output file
// based on media metadata
func (s *metadataFileHandler) handle(ctx context.Context, srcMedia mediatype.Format) error {
	visitor := mediatype.FormatWithVisitor[string](srcMedia)
	srcPath, err := visitor.Accept(ctx, visitors.NewMediaPath(ctx))
	if err != nil {
		return err
	}

	logger := ilog.FromContext(ctx).With(zap.String("sourcePath", srcPath))
	logger.Debug("Processing file...")
	outPath, err := s.getOutputFile(ctx, srcMedia)
	if err != nil {
		return err
	}

	logger = logger.With(zap.String("outPath", outPath))
	if srcPath == outPath {
		logger.Debug("Source and destination file match. Nothing to do.")
		return nil
	}

	skip, err := s.shouldSkip(ctx, srcMedia, outPath)
	if err != nil {
		return err
	}
	if skip {
		logger.Debug("Skipping moving source file...")
		return nil
	}

	if s.dryRun {
		logger.Debug("Dry run, moving file...")
		return nil
	}

	logger.Debug("Moving file...")
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	if err := os.Rename(srcPath, outPath); err != nil {
		return err
	}

	logger.Debug("Successfully moved file.")
	return nil
}

func (s *metadataFileHandler) getOutputFile(ctx context.Context, media mediatype.Format) (string, error) {
	visitor := mediatype.FormatWithVisitor[visitors.MediaMetadata](media)
	mediaMetadata, err := visitor.Accept(ctx, s.mediaMetadataVisitorFunc)
	if err != nil {
		return "", err
	}
	return mediaMetadata.OutPath, nil
}

func (s *metadataFileHandler) shouldSkip(ctx context.Context, srcMedia mediatype.Format, outPath string) (bool, error) {
	_, err := os.Stat(outPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if s.detectDuplicates {
		isDuplicate, err := s.isDuplicateImage(ctx, srcMedia, outPath)
		if err != nil {
			return false, fmt.Errorf("%w: failed to detect if images are duplicate", err)
		}
		if isDuplicate {
			ilog.FromContext(ctx).Debug("Detected duplicate image.")
			return true, nil
		}
	}

	if s.overwriteExisting && !s.dryRun {
		ilog.FromContext(ctx).Debug("Force overwrite flag set. Removing existing file at path.",
			zap.String("path", outPath))
		return false, os.Remove(outPath)
	}
	return false, fmt.Errorf("desired output filename collision")
}

func (s *metadataFileHandler) isDuplicateImage(ctx context.Context, srcMedia mediatype.Format, outPath string) (bool, error) {
	outMedia, err := mediatype.ID(outPath, s.useInputMagicSignature)
	if err != nil {
		return false, err
	}

	if !mediatype.EqualFormats(srcMedia, outMedia) {
		return false, nil
	}

	duplicateVisitorFunc, err := visitors.NewIsDuplicateMedia(ctx, srcMedia)
	if err != nil {
		return false, err
	}
	visitor := mediatype.FormatWithVisitor[bool](outMedia)
	return visitor.Accept(ctx, duplicateVisitorFunc)
}

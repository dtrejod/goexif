package mediadate

import (
	"context"

	"github.com/dtrejod/goexif/internal/ilog"
	"github.com/dtrejod/goexif/internal/mediatype"
	"github.com/dtrejod/goexif/internal/visitors"
	"go.uber.org/zap"
)

// Print logs the datetime for a provided mediafile
func Print(ctx context.Context, path string, useMagicSignature bool) error {
	media, err := mediatype.ID(path, useMagicSignature)
	if err != nil {
		return err
	}

	visitorFunc := visitors.NewMediaMetadataFilename(ctx, nil, false, false, false)
	visitor := mediatype.FormatWithVisitor[visitors.MediaMetadata](media)
	mediaMetadata, err := visitor.Accept(ctx, visitorFunc)
	if err != nil {
		return err
	}

	ilog.FromContext(ctx).Info("Found date metadata for media.",
		zap.String("sourceFile", path),
		zap.String("humanTimestamp", mediaMetadata.Timestamp.String()),
		zap.Time("unixTimestamp", mediaMetadata.Timestamp))
	return nil
}

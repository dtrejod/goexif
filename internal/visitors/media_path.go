package visitors

import (
	"context"

	"github.com/dtrejod/goexif/internal/mediatype"
)

type mediaPath struct{}

// NewMediaPath is a mediatype visitor that will get the path from a media format
func NewMediaPath(ctx context.Context) mediatype.VisitorFunc[string] {
	return &mediaPath{}
}

func (m *mediaPath) VisitJPEG(_ context.Context, media mediatype.JPEG) (string, error) {
	return media.Path, nil
}

func (m *mediaPath) VisitPNG(_ context.Context, media mediatype.PNG) (string, error) {
	return media.Path, nil
}

func (m *mediaPath) VisitHEIF(_ context.Context, media mediatype.HEIF) (string, error) {
	return media.Path, nil
}

func (m *mediaPath) VisitTIFF(_ context.Context, media mediatype.TIFF) (string, error) {
	return media.Path, nil
}

func (m *mediaPath) VisitQTFF(_ context.Context, media mediatype.QTFF) (string, error) {
	return media.Path, nil
}

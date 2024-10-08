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

func (m *mediaPath) VisitMP4(_ context.Context, media mediatype.MP4) (string, error) {
	return media.Path, nil
}

func (m *mediaPath) VisitAVI(_ context.Context, media mediatype.AVI) (string, error) {
	return media.Path, nil
}

func (m *mediaPath) Visit3PG(_ context.Context, media mediatype.GPP) (string, error) {
	return media.Path, nil
}

func (m *mediaPath) Visit3G2(_ context.Context, media mediatype.GPP2) (string, error) {
	return media.Path, nil
}

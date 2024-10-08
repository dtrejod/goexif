package visitors

import (
	"context"

	"github.com/dtrejod/goexif/internal/mediatype"
)

type mediaExt struct{}

// NewMediaExtAliases is a mediatype visitor that will get all known media extension aliases
func NewMediaExtAliases(ctx context.Context) mediatype.VisitorFunc[map[string]struct{}] {
	return &mediaExt{}
}

func (m *mediaExt) VisitJPEG(_ context.Context, media mediatype.JPEG) (map[string]struct{}, error) {
	return media.Aliases(), nil
}

func (m *mediaExt) VisitPNG(_ context.Context, media mediatype.PNG) (map[string]struct{}, error) {
	return media.Aliases(), nil
}

func (m *mediaExt) VisitHEIF(_ context.Context, media mediatype.HEIF) (map[string]struct{}, error) {
	return media.Aliases(), nil
}

func (m *mediaExt) VisitTIFF(_ context.Context, media mediatype.TIFF) (map[string]struct{}, error) {
	return media.Aliases(), nil
}

func (m *mediaExt) VisitQTFF(_ context.Context, media mediatype.QTFF) (map[string]struct{}, error) {
	return media.Aliases(), nil
}

func (m *mediaExt) VisitMP4(_ context.Context, media mediatype.MP4) (map[string]struct{}, error) {
	return media.Aliases(), nil
}

func (m *mediaExt) VisitAVI(_ context.Context, media mediatype.AVI) (map[string]struct{}, error) {
	return media.Aliases(), nil
}

func (m *mediaExt) Visit3PG(_ context.Context, media mediatype.GPP) (map[string]struct{}, error) {
	return media.Aliases(), nil
}

func (m *mediaExt) Visit3G2(_ context.Context, media mediatype.GPP2) (map[string]struct{}, error) {
	return media.Aliases(), nil
}

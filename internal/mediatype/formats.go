package mediatype

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
)

// MediaType represents a common interface for all media formats
type MediaType interface {
	// String implements Stringer interface
	String() string
	// Ext returns the file extension
	Ext() string
	// Aliases returns known file type aliases for this media type
	Aliases() map[string]struct{}
}

// Format is a container for any known media type
type Format struct {
	media MediaType
}

// AllKnownMediaTypes is a list of all known media types
var AllKnownMediaTypes = []MediaType{
	JPEG{},
	PNG{},
	HEIF{},
	TIFF{},
	QTFF{},
	MP4{},
	AVI{},
	GPP{},
	GPP2{},
}

// NewFormat returns a new Format instance based on the file extension. If useSignature is true, then the existing file
// extension is ignored and we use the file magic signature instead
// REF: https://en.wikipedia.org/wiki/File_format#Magic_number
func NewFormat(path string, useSignature bool) (Format, error) {
	ext := strings.TrimPrefix(filepath.Ext(strings.ToLower(path)), ".")
	if useSignature {
		t, err := filetype.MatchFile(path)
		if err != nil {
			return Format{}, err
		}
		ext = t.Extension
	}
	switch {
	case contains(JPEG{}.Aliases(), ext):
		return Format{media: JPEG{Path: path}}, nil
	case contains(PNG{}.Aliases(), ext):
		return Format{media: PNG{Path: path}}, nil
	case contains(HEIF{}.Aliases(), ext):
		return Format{media: HEIF{Path: path}}, nil
	case contains(TIFF{}.Aliases(), ext):
		return Format{media: TIFF{Path: path}}, nil
	case contains(QTFF{}.Aliases(), ext):
		return Format{media: QTFF{Path: path}}, nil
	case contains(MP4{}.Aliases(), ext):
		return Format{media: MP4{Path: path}}, nil
	case contains(AVI{}.Aliases(), ext):
		return Format{media: AVI{Path: path}}, nil
	case contains(GPP{}.Aliases(), ext):
		return Format{media: GPP{Path: path}}, nil
	case contains(GPP2{}.Aliases(), ext):
		return Format{media: GPP2{Path: path}}, nil
	default:
		return Format{media: Unknown{}}, nil
	}
}

func contains(toMatch map[string]struct{}, s string) bool {
	_, ok := toMatch[s]
	return ok
}

// FormatWithVisitor is a generic Format union type visitor
type FormatWithVisitor[T any] Format

// VisitorFunc implements a Visitor type that handles all known Format types
type VisitorFunc[T any] interface {
	VisitJPEG(context.Context, JPEG) (T, error)
	VisitPNG(context.Context, PNG) (T, error)
	VisitHEIF(context.Context, HEIF) (T, error)
	VisitTIFF(context.Context, TIFF) (T, error)
	VisitQTFF(context.Context, QTFF) (T, error)
	VisitMP4(context.Context, MP4) (T, error)
	VisitAVI(context.Context, AVI) (T, error)
	Visit3PG(context.Context, GPP) (T, error)
	Visit3G2(context.Context, GPP2) (T, error)
}

// Accept visits the current media type using the visitor pattern
func (f *FormatWithVisitor[T]) Accept(ctx context.Context, v VisitorFunc[T]) (T, error) {
	switch f.media.(type) {
	case JPEG:
		return v.VisitJPEG(ctx, f.media.(JPEG))
	case PNG:
		return v.VisitPNG(ctx, f.media.(PNG))
	case HEIF:
		return v.VisitHEIF(ctx, f.media.(HEIF))
	case TIFF:
		return v.VisitTIFF(ctx, f.media.(TIFF))
	case QTFF:
		return v.VisitQTFF(ctx, f.media.(QTFF))
	case MP4:
		return v.VisitMP4(ctx, f.media.(MP4))
	case AVI:
		return v.VisitAVI(ctx, f.media.(AVI))
	case GPP:
		return v.Visit3PG(ctx, f.media.(GPP))
	case GPP2:
		return v.Visit3G2(ctx, f.media.(GPP2))
	case Unknown:
	default:
	}
	return *new(T), fmt.Errorf("unknown media type")
}

// EqualFormats returns true if two Formats are of the same media type
func EqualFormats(a, b Format) bool {
	// Use type assertions to compare the types of media
	switch a.media.(type) {
	case JPEG:
		_, ok := b.media.(JPEG)
		return ok
	case PNG:
		_, ok := b.media.(PNG)
		return ok
	case HEIF:
		_, ok := b.media.(HEIF)
		return ok
	case TIFF:
		_, ok := b.media.(TIFF)
		return ok
	case QTFF:
		_, ok := b.media.(QTFF)
		return ok
	case MP4:
		_, ok := b.media.(MP4)
		return ok
	case AVI:
		_, ok := b.media.(AVI)
		return ok
	case GPP:
		_, ok := b.media.(GPP)
		return ok
	case GPP2:
		_, ok := b.media.(GPP2)
		return ok
	case Unknown:
		_, ok := b.media.(Unknown)
		return ok
	default:
		return false
	}
}

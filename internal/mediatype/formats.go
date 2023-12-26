package mediatype

import (
	"context"
	"fmt"
)

// Format is union type of all known mediatype formats
type Format struct {
	t    mediaFormat
	jpeg JPEG
	png  PNG
	heif HEIF
}

type mediaFormat string

const (
	jpegMediaFormat mediaFormat = "jpeg"
	pngMediaFormat  mediaFormat = "png"
	heifMediaFormat mediaFormat = "heif"
)

// NewJPEGFormat returns a new JPEG format
func NewJPEGFormat(j JPEG) Format {
	return Format{t: jpegMediaFormat, jpeg: j}
}

// NewPNGFormat returns a new PNG format
func NewPNGFormat(p PNG) Format {
	return Format{t: pngMediaFormat, png: p}
}

// NewHEIFFormat returns a new HEIF format
func NewHEIFFormat(h HEIF) Format {
	return Format{t: heifMediaFormat, heif: h}
}

// FormatWithVisitor is a generic Format union type visitor
type FormatWithVisitor[T any] Format

// Accept visits all known Format types
func (f *FormatWithVisitor[T]) Accept(ctx context.Context, v VisitorFunc[T]) (T, error) {
	switch f.t {
	case jpegMediaFormat:
		return v.VisitJPEG(ctx, f.jpeg)
	case pngMediaFormat:
		return v.VisitPNG(ctx, f.png)
	case heifMediaFormat:
		return v.VisitHEIF(ctx, f.heif)
	default:
		return *new(T), fmt.Errorf("unknown media type: %s", string(f.t))
	}
}

// VisitorFunc implements a Visitor type that handles all known Format types
type VisitorFunc[T any] interface {
	VisitJPEG(context.Context, JPEG) (T, error)
	VisitPNG(context.Context, PNG) (T, error)
	VisitHEIF(context.Context, HEIF) (T, error)
}

// EqualFormats returns true if two Formats are the same type
func EqualFormats(a, b Format) bool {
	return a.t == b.t
}

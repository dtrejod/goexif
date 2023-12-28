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
	tiff TIFF
	mov  QTFF
}

type mediaFormat string

const (
	jpegMediaFormat mediaFormat = "jpeg"
	pngMediaFormat  mediaFormat = "png"
	heifMediaFormat mediaFormat = "heif"
	tiffMediaFormat mediaFormat = "tiff"
	qtffMediaFormat mediaFormat = "qtff"
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

// NewTIFFFormat returns a new TIFF format
func NewTIFFFormat(h TIFF) Format {
	return Format{t: tiffMediaFormat, tiff: h}
}

// NewQTFFFormat returns a new Quicktime format
func NewQTFFFormat(h QTFF) Format {
	return Format{t: qtffMediaFormat, mov: h}
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
	case tiffMediaFormat:
		return v.VisitTIFF(ctx, f.tiff)
	case qtffMediaFormat:
		return v.VisitQTFF(ctx, f.mov)
	default:
		return *new(T), fmt.Errorf("unknown media type")
	}
}

// VisitorFunc implements a Visitor type that handles all known Format types
type VisitorFunc[T any] interface {
	VisitJPEG(context.Context, JPEG) (T, error)
	VisitPNG(context.Context, PNG) (T, error)
	VisitHEIF(context.Context, HEIF) (T, error)
	VisitTIFF(context.Context, TIFF) (T, error)
	VisitQTFF(context.Context, QTFF) (T, error)
}

// EqualFormats returns true if two Formats are the same type
func EqualFormats(a, b Format) bool {
	return a.t == b.t
}

package mediatype

import (
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
)

// ID identifies media into a known mediatype format. If useSignature is true,
// then the existing file extension is ignored and we use the file magic
// signature instead
// REF: https://en.wikipedia.org/wiki/File_format#Magic_number
func ID(path string, useSignature bool) (Format, error) {
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
		return NewJPEGFormat(JPEG{Path: path}), nil
	case contains(PNG{}.Aliases(), ext):
		return NewPNGFormat(PNG{Path: path}), nil
	case contains(HEIF{}.Aliases(), ext):
		return NewHEIFFormat(HEIF{Path: path}), nil
	case contains(TIFF{}.Aliases(), ext):
		return NewTIFFFormat(TIFF{Path: path}), nil
	case contains(QTFF{}.Aliases(), ext):
		return NewQTFFFormat(QTFF{Path: path}), nil
	case contains(MP4{}.Aliases(), ext):
		return NewMP4Format(MP4{Path: path}), nil
	case contains(AVI{}.Aliases(), ext):
		return NewAVIFormat(AVI{Path: path}), nil
	default:
		return Format{}, nil
	}
}

func contains(toMatch map[string]struct{}, s string) bool {
	_, ok := toMatch[s]
	return ok
}

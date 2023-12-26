package mediatype

// PNG identifies PNG media.
// EXIF extension was adopted for PNG in 2017
// http://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html#C.eXIf
type PNG struct {
	Path string
}

// String implements Stringer interface
func (t PNG) String() string {
	return "png"
}

// Ext returns the file extension
func (t PNG) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t PNG) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
	}
}

package mediatype

// TIFF indetifies TIFF media
type TIFF struct {
	Path string
}

// String implements Stringer interface
func (t TIFF) String() string {
	return "tiff"
}

// Ext returns the file extension
func (t TIFF) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t TIFF) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
	}
}

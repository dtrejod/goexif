package mediatype

// JPEG indetifies JPEG media
// REF: https://en.wikipedia.org/wiki/JPEG
type JPEG struct {
	Path string
}

// String implements Stringer interface
func (t JPEG) String() string {
	return "jpg"
}

// Ext returns the file extension
func (t JPEG) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t JPEG) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
		"jpeg":     {},
		"jpe":      {},
		"jif":      {},
		"jfif":     {},
		"jfi":      {},
	}
}

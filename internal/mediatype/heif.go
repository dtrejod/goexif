package mediatype

// HEIF identifies HEIF media
type HEIF struct {
	Path string
}

// String implements Stringer interface
func (t HEIF) String() string {
	return "heif"
}

// Ext returns the file extension
func (t HEIF) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t HEIF) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
	}
}

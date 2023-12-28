package mediatype

// MP4 indetifies Quicktime media
type MP4 struct {
	Path string
}

// String implements Stringer interface
func (t MP4) String() string {
	return "mp4"
}

// Ext returns the file extension
func (t MP4) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t MP4) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
	}
}

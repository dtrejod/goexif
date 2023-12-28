package mediatype

// MP4 indetifies Quicktime media
// Ref: https://en.wikipedia.org/wiki/MP4_file_format
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
		"m4a":      {},
		"m4p":      {},
		"m4b":      {},
		"m4r":      {},
		"m4v":      {},
	}
}

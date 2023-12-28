package mediatype

// QTFF indetifies Quicktime media
// Ref: https://en.wikipedia.org/wiki/QuickTime_File_Format
type QTFF struct {
	Path string
}

// String implements Stringer interface
func (t QTFF) String() string {
	return "mov"
}

// Ext returns the file extension
func (t QTFF) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t QTFF) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
		"movie":    {},
		"qt":       {},
	}
}

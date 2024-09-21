package mediatype

// Unknown identifies an unknown media type
type Unknown struct {
	Path string
}

// String implements Stringer interface
func (t Unknown) String() string {
	return ""
}

// Ext returns the file extension
func (t Unknown) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t Unknown) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
	}
}

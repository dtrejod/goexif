package mediatype

// AVI indetifies AVI media
// REF: https://en.wikipedia.org/wiki/Audio_Video_Interleave
type AVI struct {
	Path string
}

// String implements Stringer interface
func (t AVI) String() string {
	return "avi"
}

// Ext returns the file extension
func (t AVI) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t AVI) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
	}
}

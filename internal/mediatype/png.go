package mediatype

// PNG identifies PNG media.
// REF: https://en.wikipedia.org/wiki/PNG
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

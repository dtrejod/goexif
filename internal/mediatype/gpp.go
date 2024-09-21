package mediatype

// GPP indetifies 3GP media
// REF: https://en.wikipedia.org/wiki/3GP_and_3G2
type GPP struct {
	Path string
}

// String implements Stringer interface
func (t GPP) String() string {
	return "3gp"
}

// Ext returns the file extension
func (t GPP) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t GPP) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
		"3gpp":     {},
	}
}

// GPP2 indetifies 3G2 media
// REF: https://en.wikipedia.org/wiki/3GP_and_3G2
type GPP2 struct {
	Path string
}

// String implements Stringer interface
func (t GPP2) String() string {
	return "3g2"
}

// Ext returns the file extension
func (t GPP2) Ext() string {
	return "." + t.String()
}

// Aliases returns known file type aliases for this media type
func (t GPP2) Aliases() map[string]struct{} {
	return map[string]struct{}{
		t.String(): {},
		"3gp2":     {},
		"3gpp2":    {},
	}
}

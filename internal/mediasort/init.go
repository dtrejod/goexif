package mediasort

import "github.com/dtrejod/goexif/internal/mediatype"

var (
	// DefaultFileTypes are the default media types handled by the sorter if none are specified.
	// NOTE: This list is generated on pkg init
	DefaultFileTypes = []string{}
)

// init initializes the FileTypes variable
func init() {
	for _, media := range mediatype.AllKnownMediaTypes {
		DefaultFileTypes = append(DefaultFileTypes, media.String())
	}
}

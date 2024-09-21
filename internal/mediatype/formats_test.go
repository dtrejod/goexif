package mediatype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualFormats(t *testing.T) {
	// Test that all known media types are equal to themselves
	for _, mt := range AllKnownMediaTypes {
		f1 := Format{media: mt}
		f2 := Format{media: mt}
		assert.True(t, EqualFormats(f1, f2))
	}
}

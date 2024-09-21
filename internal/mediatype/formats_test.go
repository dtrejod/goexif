package mediatype

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormat(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Test that NewFormat() returns a Format with the correct media type
	for _, mt := range AllKnownMediaTypes {
		// Test with a file extension
		tmpFile, err := os.CreateTemp(tmpDir, fmt.Sprintf("*.%s", mt.Ext()))
		defer tmpFile.Close()
		require.NoError(t, err)

		f, err := NewFormat(tmpFile.Name(), false)
		assert.NoError(t, err)
		assert.True(t, EqualFormats(f, Format{media: mt}))
	}
}

func TestEqualFormats(t *testing.T) {
	// Test that all known media types are equal to themselves
	for _, mt := range AllKnownMediaTypes {
		f1 := Format{media: mt}
		f2 := Format{media: mt}
		assert.True(t, EqualFormats(f1, f2))
	}
}

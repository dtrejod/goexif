package visitors

import (
	"context"
	"testing"
	"time"

	"github.com/dtrejod/goexif/internal/mediatype"
	"github.com/stretchr/testify/assert"
)

func TestMetadataFilename(t *testing.T) {
	ctx := context.Background()

	t.Run("with default config", func(t *testing.T) {
		expected := MediaMetadata{
			OutPath:   "2000/01/01/white.png",
			Timestamp: time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		}
		srcMedia, err := mediatype.ID("./testdata/white.png", false)
		assert.NoError(t, err)

		visitorFunc := NewMediaMetadataFilename(ctx, toPtr("."), false, false, false)
		visitor := mediatype.FormatWithVisitor[MediaMetadata](srcMedia)
		actual, err := visitor.Accept(ctx, visitorFunc)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	// TODO(dtrejo): Fix mod time test. Git doesn't preserve the mod time of files making this test more involved
	//t.Run("with last modified timestamp", func(t *testing.T) {
	//	expected := "2001/01/01/noexif.png"
	//	srcMedia, err := mediatype.ID("./testdata/noexif.png", false)
	//	assert.NoError(t, err)

	//	visitorFunc := NewMediaMetadataFilename(ctx, toPtr("."), true, false, false)
	//	visitor := mediatype.FormatWithVisitor[string](srcMedia)
	//	actual, err := visitor.Accept(ctx, visitorFunc)
	//	assert.NoError(t, err)

	//	assert.Equal(t, expected, actual)
	//})

	t.Run("with timestamp as filename", func(t *testing.T) {
		expected := MediaMetadata{
			OutPath:   "2000/01/01/946684800.png",
			Timestamp: time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		}
		srcMedia, err := mediatype.ID("./testdata/white.png", false)
		assert.NoError(t, err)

		visitorFunc := NewMediaMetadataFilename(ctx, toPtr("."), false, true, false)
		visitor := mediatype.FormatWithVisitor[MediaMetadata](srcMedia)
		actual, err := visitor.Accept(ctx, visitorFunc)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("with clean file extension", func(t *testing.T) {
		expected := MediaMetadata{
			OutPath:   "2000/01/01/ispng.png",
			Timestamp: time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		}
		srcMedia, err := mediatype.ID("./testdata/ispng.jpg", true)
		assert.NoError(t, err)

		visitorFunc := NewMediaMetadataFilename(ctx, toPtr("."), false, false, true)
		visitor := mediatype.FormatWithVisitor[MediaMetadata](srcMedia)
		actual, err := visitor.Accept(ctx, visitorFunc)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}

func toPtr[T any](v T) *T {
	return &v
}

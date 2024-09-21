package visitors

import (
	"context"
	"testing"

	"github.com/dtrejod/goexif/internal/mediatype"
	"github.com/stretchr/testify/assert"
)

func TestIsDuplicateDectector(t *testing.T) {
	ctx := context.Background()

	t.Run("compare white png are identical", func(t *testing.T) {
		white, err := mediatype.NewFormat("./testdata/white.png", false)
		assert.NoError(t, err)

		visitorFunc, err := NewIsDuplicateMedia(ctx, white)
		assert.NoError(t, err)
		visitor := mediatype.FormatWithVisitor[bool](white)
		actual, err := visitor.Accept(ctx, visitorFunc)
		assert.NoError(t, err)

		assert.True(t, actual)
	})

	t.Run("compare white and black png are identical", func(t *testing.T) {
		// phash ignores color changes so white and black are considered the same
		white, err := mediatype.NewFormat("./testdata/white.png", false)
		assert.NoError(t, err)
		black, err := mediatype.NewFormat("./testdata/black.png", false)
		assert.NoError(t, err)

		visitorFunc, err := NewIsDuplicateMedia(ctx, white)
		assert.NoError(t, err)
		visitor := mediatype.FormatWithVisitor[bool](black)
		actual, err := visitor.Accept(ctx, visitorFunc)
		assert.NoError(t, err)

		assert.True(t, actual)
	})

	t.Run("compare identical monas are duplicates", func(t *testing.T) {
		mona, err := mediatype.NewFormat("./testdata/mona/mona.jpg", false)
		assert.NoError(t, err)

		visitorFunc, err := NewIsDuplicateMedia(ctx, mona)
		assert.NoError(t, err)
		visitor := mediatype.FormatWithVisitor[bool](mona)
		actual, err := visitor.Accept(ctx, visitorFunc)
		assert.NoError(t, err)

		assert.True(t, actual)
	})

	t.Run("compare original and dark mona are duplicates", func(t *testing.T) {
		mona, err := mediatype.NewFormat("./testdata/mona/mona.jpg", false)
		assert.NoError(t, err)
		dark, err := mediatype.NewFormat("./testdata/mona/dark.jpg", false)
		assert.NoError(t, err)

		visitorFunc, err := NewIsDuplicateMedia(ctx, mona)
		assert.NoError(t, err)
		visitor := mediatype.FormatWithVisitor[bool](dark)
		actual, err := visitor.Accept(ctx, visitorFunc)
		assert.NoError(t, err)

		assert.True(t, actual)
	})

	t.Run("compare original and color mona are duplicates", func(t *testing.T) {
		mona, err := mediatype.NewFormat("./testdata/mona/mona.jpg", false)
		assert.NoError(t, err)
		color, err := mediatype.NewFormat("./testdata/mona/color.jpg", false)
		assert.NoError(t, err)

		visitorFunc, err := NewIsDuplicateMedia(ctx, mona)
		assert.NoError(t, err)
		visitor := mediatype.FormatWithVisitor[bool](color)
		actual, err := visitor.Accept(ctx, visitorFunc)
		assert.NoError(t, err)

		assert.True(t, actual)
	})

	t.Run("compare original and raphael's mona are different", func(t *testing.T) {
		mona, err := mediatype.NewFormat("./testdata/mona/mona.jpg", false)
		assert.NoError(t, err)
		raphael, err := mediatype.NewFormat("./testdata/mona/raphael.jpg", false)
		assert.NoError(t, err)

		visitorFunc, err := NewIsDuplicateMedia(ctx, mona)
		assert.NoError(t, err)
		visitor := mediatype.FormatWithVisitor[bool](raphael)
		actual, err := visitor.Accept(ctx, visitorFunc)
		assert.NoError(t, err)

		assert.False(t, actual)
	})
}

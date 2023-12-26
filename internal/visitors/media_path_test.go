package visitors

import (
	"context"
	"testing"

	"github.com/dtrejod/goexif/internal/mediatype"
	"github.com/stretchr/testify/assert"
)

func TestMediaPath(t *testing.T) {
	ctx := context.Background()
	expected := "./testdata/white.png"
	srcMedia, err := mediatype.ID("./testdata/white.png", false)
	assert.NoError(t, err)

	visitorFunc := NewMediaPath(ctx)
	visitor := mediatype.FormatWithVisitor[string](srcMedia)
	actual, err := visitor.Accept(ctx, visitorFunc)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

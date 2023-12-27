package visitors

import (
	"context"
	"fmt"
	"image"
	"os"

	// jpeg import for side effect of decoding jpeg images
	_ "image/jpeg"
	// png import for side effect of decoding png images
	_ "image/png"
	// tiff import for side effect of decoding tiff images
	_ "golang.org/x/image/tiff"

	"github.com/corona10/goimagehash"
	"github.com/dtrejod/goexif/internal/ilog"
	"github.com/dtrejod/goexif/internal/mediatype"
	"go.uber.org/zap"
)

const (
	// pHashDistanceUpperBound is the upper bound distance between 2 image
	// perception hashes that determines when the images are same
	// A hash distance > 50 is a typical threshold for different images.
	pHashDistanceUpperBound = 25
)

type mediaCompare struct {
	srcPath string
}

// NewIsDuplicateMedia is a mediatype visitor that will determine whether 2 media formats are the same
func NewIsDuplicateMedia(ctx context.Context, srcMedia mediatype.Format) (mediatype.VisitorFunc[bool], error) {
	visitor := mediatype.FormatWithVisitor[string](srcMedia)
	srcPath, err := visitor.Accept(ctx, NewMediaPath(ctx))
	if err != nil {
		return nil, err
	}
	return &mediaCompare{
		srcPath: srcPath,
	}, nil
}

func (m *mediaCompare) VisitJPEG(ctx context.Context, outMedia mediatype.JPEG) (bool, error) {
	return compareUsingPHash(ctx, m.srcPath, outMedia.Path)
}

func (m *mediaCompare) VisitPNG(ctx context.Context, outMedia mediatype.PNG) (bool, error) {
	return compareUsingPHash(ctx, m.srcPath, outMedia.Path)
}

func (m *mediaCompare) VisitHEIF(ctx context.Context, outMedia mediatype.HEIF) (bool, error) {
	return false, fmt.Errorf("checking for duplicate is not supported for heif media")
}

func (m *mediaCompare) VisitTIFF(ctx context.Context, outMedia mediatype.TIFF) (bool, error) {
	return compareUsingPHash(ctx, m.srcPath, outMedia.Path)
}

func compareUsingPHash(ctx context.Context, src, dest string) (bool, error) {
	logger := ilog.FromContext(ctx).With(
		zap.String("sourcePath", src),
		zap.String("destinationPath", dest))

	hashA, err := getImagePerceptionHash(src)
	if err != nil {
		return false, err
	}

	hashB, err := getImagePerceptionHash(dest)
	if err != nil {
		return false, err
	}

	distance, err := hashA.Distance(hashB)
	if err != nil {
		return false, fmt.Errorf("%w: failed to compare images for duplicates", err)
	}

	logger.Debug("Comparing images...",
		zap.String("hashA", hashA.ToString()),
		zap.String("hashB", hashB.ToString()),
		zap.Int("distance", distance))

	return distance < pHashDistanceUpperBound, nil
}

// getImagePerceptionHash returns the peception hash of the image. The phash is
// taken first performing a discrete cosine transform on the image. Then it
// compares each pixel to it's average. If it is larger then output 1 else 0
// otherwise.
// Since a Phash operates in the frequency domain, it should be more tolerant
// to color shifts, size scaling, watermarks, and compression between 2 images.
func getImagePerceptionHash(path string) (*goimagehash.ImageHash, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to open image", err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to decode image", err)
	}

	hashA, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get perception hash image", err)
	}
	return hashA, nil
}
